package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/metaclips/LetsTalk/values"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

func init() {
	var m = webrtc.MediaEngine{}

	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	classSessions.api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
}

var (
	// Create a MediaEngine object to configure the supported codec
	classSessions = classSessionPeerConnections{
		publisherVideoTracks:  make(map[string]*webrtc.Track),
		publisherTrackMutexes: &sync.RWMutex{},

		audioTracks:       make(map[string][]*webrtc.Track),
		audioTrackSender:  make(map[*webrtc.Track][]rtpSenderData),
		audioTrackMutexes: &sync.RWMutex{},

		peerConnection:        make(map[string]*webrtc.PeerConnection),
		peerConnectionMutexes: &sync.RWMutex{},

		connectedUsers:      make(map[string][]string),
		connectedUsersMutex: &sync.RWMutex{},
	}
)

func (s *classSessionPeerConnections) startClassSession(msg []byte, user string) {
	sessionID := uuid.New().String()
	sdp := sdpConstruct{}

	if err := json.Unmarshal(msg, &sdp); err != nil {
		// Send back a CreateSessionError indicating user already in a session.
		onSessionError(user, "Unable to retrieve class session details.")
		log.Println("error while unmarshaling start session sdp", err)

		return
	}

	// A single user might login using multiple devices. We close recent peerconnection if there's one.
	s.peerConnectionMutexes.RLock()
	if s.peerConnection[sdp.UserID] != nil {
		log.Println(user, "already in a session")
		onSessionError(user, "You are already in another session.")
		s.peerConnectionMutexes.RUnlock()
		return
	}
	s.peerConnectionMutexes.RUnlock()

	peerConnection, err := classSessions.api.NewPeerConnection(values.PeerConnectionConfig)
	if err != nil {
		log.Println("unable to create a peerconnection", err)
		onSessionError(user, "Unable to create peerconnection")
		return
	}

	// Add peerconnection to map.
	s.peerConnectionMutexes.Lock()
	s.peerConnection[sdp.UserID] = peerConnection
	s.peerConnectionMutexes.Unlock()

	sdp.ClassSessionID = sessionID

	s.connectedUsersMutex.Lock()
	s.connectedUsers[sessionID] = []string{sdp.UserID}
	s.connectedUsersMutex.Unlock()

	var videoAudioWriter *webmWriter
	if values.Config.EnableClassSessionRecord {
		videoAudioWriter = newWebmWriter(sessionID + ".webm")
	}

	peerConnection.OnConnectionStateChange(func(cc webrtc.PeerConnectionState) {
		fmt.Printf("PeerConnection State has changed %s \n", cc.String())
		if cc == webrtc.PeerConnectionStateFailed {
			// Since this is the publisher, all video and audio tracks related to the session
			// should be cleared and all peer connections closed.
			videoAudioWriter.close()

			go func() {
				if !values.Config.EnableClassSessionRecord {
					return
				}

				// If token is not provided and file upload is set to true, file is uploaded to DB.
				if values.Config.DropboxToken == "" {
					videoAudioWriter.uploadToDB()
					return
				}

				link, err := videoAudioWriter.getVideoFileSharableLink()
				if err != nil {
					log.Println("unable to generate sharable link", err)
					return
				}

				if len(link) > 0 {
					link = link[:len(link)-1] + "1"
				}

				message := Message{
					RoomID:      sdp.RoomID,
					Name:        sdp.AuthorName,
					Message:     link,
					UserID:      sdp.UserID,
					Type:        "classSessionLink",
					MessageType: "ClassSessionLink",
					FileHash:    sdp.ClassSessionID,
				}

				roomUsers, err := message.saveMessageContent()
				if err != nil {
					log.Println("unable to save class session link to db", err)
				}

				content, err := json.Marshal(message)
				if err != nil {
					log.Println("unable to marshal json while sending link", err)
					return
				}

				for _, roomUser := range roomUsers {
					HubConstruct.sendMessage(content, roomUser)
				}
			}()

			s.peerConnectionMutexes.RLock()
			s.connectedUsersMutex.Lock()

			for _, user := range s.connectedUsers[sessionID] {
				closePeerConnection(s.peerConnection[user])

				delete(s.peerConnection, user)
			}

			delete(s.connectedUsers, sessionID)

			s.connectedUsersMutex.Unlock()
			s.peerConnectionMutexes.RUnlock()

			s.audioTrackMutexes.Lock()
			delete(s.audioTracks, sessionID)
			s.audioTrackMutexes.Unlock()

			s.publisherTrackMutexes.Lock()
			delete(s.publisherVideoTracks, sessionID)
			s.publisherTrackMutexes.Unlock()
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnSignalingStateChange(func(cc webrtc.SignalingState) {
		fmt.Println("Session singaling", cc.String())
	})

	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Publisher <-> Server is to receive both audio and video packets.
		// Packets are to be broadcasted to other users on the session.
		// For video, resolution is in 480px.
		if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP8 {
			log.Println("VP8 track is being called")

			videoTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sessionID, sdp.UserID)
			if err != nil {
				log.Println("unable to generate start session video track", err)
				onSessionError(user, "Unable to start video track.")

				// Return back a class session creation error back to client.
				peerConnection.Close()
				return
			}

			s.publisherTrackMutexes.Lock()
			s.publisherVideoTracks[sessionID] = videoTrack
			s.publisherTrackMutexes.Unlock()

			// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
			go func() {
				ticker := time.NewTicker(values.RtcpPLIInterval)
				for range ticker.C {
					err := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: videoTrack.SSRC()}})
					if err != nil {
						break
					}
				}
			}()

			for {
				rtp, err := remoteTrack.ReadRTP()
				if err != nil {
					log.Println("Publisher video track errored, exiting now.", err)
					break
				}

				err = videoTrack.WriteRTP(rtp)
				if err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Println("Publisher video packet writed break", err)
					break
				}

				if values.Config.EnableClassSessionRecord {
					videoAudioWriter.pushVP8(rtp)
				}
			}

			log.Println("Publisher video track exited")

		} else if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeOpus {
			log.Println("OPUS track called")

			audioTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sessionID, sdp.UserID)
			if err != nil {
				log.Println("unable to start audio track", err)
				// Return back a class session creation error back to client.
				// Also, users might decide to disable video/audio on start.
				onSessionError(user, "Unable to start audio track.")
				peerConnection.Close()
				return
			}

			s.audioTrackMutexes.Lock()
			s.audioTracks[sessionID] = []*webrtc.Track{audioTrack}
			s.audioTrackMutexes.Unlock()

			for {
				rtp, err := remoteTrack.ReadRTP()
				if err != nil {
					log.Println("Publisher audio track errored, exiting now.", err)
					break
				}

				err = audioTrack.WriteRTP(rtp)
				if err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Println("Publisher video packet writed break", err)
					break
				}

				if values.Config.EnableClassSessionRecord {
					videoAudioWriter.pushOpus(rtp)
				}
			}

			log.Println("Publisher audio track exited")

		} else {
			log.Println("unsupported track is being played. Video writer might not work", remoteTrack.PayloadType())
		}
	})

	sdp.peerConnection = peerConnection
	if err = sdp.negotiate(); err != nil {
		closePeerConnection(peerConnection)
		log.Println("error while negotiating on start class session", err)
		onSessionError(user, "Unable to negotiate.")

		return
	}

	// Broadcast class session to room.
	sdp.MsgType = "ClassSession"
	sdp.AuthorName = values.MapEmailToName[sdp.UserID]

	jsonContent, err := json.Marshal(sdp)
	if err != nil {
		closePeerConnection(peerConnection)
		onSessionError(user, "Unable to send class session to room.")

		return
	}

	// Send back answer SDP to client and also class notification to all users in room.
	roomUsers, err := Message{
		RoomID:   sdp.RoomID,
		Name:     sdp.AuthorName,
		UserID:   sdp.UserID,
		Type:     "classSession",
		FileHash: sdp.ClassSessionID,
	}.saveMessageContent()

	if err != nil {
		closePeerConnection(peerConnection)
		onSessionError(user, "Unable to send class session to room.")

		return
	}

	for _, user := range roomUsers {
		HubConstruct.sendMessage(jsonContent, user)
	}
}

func (s *classSessionPeerConnections) joinClassSession(msg []byte, user string) {
	sdp := sdpConstruct{}
	if err := json.Unmarshal(msg, &sdp); err != nil {
		log.Println("unable to marshal json on join session", err)
		onSessionError(user, "Unable to retrieve class session details.")

		return
	}

	// A single user might login using multiple devices. We close recent peerconnection if there's one.
	s.peerConnectionMutexes.RLock()
	if s.peerConnection[sdp.UserID] != nil {
		log.Println(user, "already in a session")
		onSessionError(user, "You are already in another session.")
		s.peerConnectionMutexes.RUnlock()

		return
	}
	s.peerConnectionMutexes.RUnlock()

	peerConnection, err := classSessions.api.NewPeerConnection(values.PeerConnectionConfig)
	if err != nil {
		log.Println("unable to create a peerconnection", err)
		onSessionError(user, "Unable to create peerconnection")

		peerConnection.Close()
		return
	}

	peerConnection.OnConnectionStateChange(func(cc webrtc.PeerConnectionState) {
		fmt.Printf("PeerConnection State has changed for joined class %s \n", cc.String())

		if cc == webrtc.PeerConnectionStateFailed {
			// Remove user from connected users list.
			s.connectedUsersMutex.Lock()
			connectedUsers := s.connectedUsers[sdp.ClassSessionID]
			for i := range connectedUsers {
				if connectedUsers[i] == sdp.UserID {
					if len(connectedUsers) > i+1 {
						connectedUsers = append(connectedUsers[:i], connectedUsers[i+1:]...)
					} else {
						connectedUsers = connectedUsers[:i]
					}

					break
				}
			}

			s.connectedUsers[sdp.ClassSessionID] = connectedUsers
			s.connectedUsersMutex.Unlock()

			// Remove audio track from other users and make a renegotiation offer.
			s.peerConnectionMutexes.Lock()
			s.audioTrackMutexes.Lock()

			closePeerConnection(peerConnection)
			delete(s.peerConnection, sdp.UserID)

			audioTracks := s.audioTracks[sdp.ClassSessionID]

			for i, audioTrack := range audioTracks {
				if audioTrack != nil && audioTrack.Label() == sdp.UserID {
					if len(audioTracks) > i+1 {
						audioTracks = append(audioTracks[:i], audioTracks[i+1:]...)
					} else {
						audioTracks = audioTracks[:i]
					}

					// Remove current subscriber track all registered tracks associated in session.
					for _, rtpSenderDetails := range s.audioTrackSender[audioTrack] {
						if pc := s.peerConnection[rtpSenderDetails.userID]; pc != nil {
							if err := pc.RemoveTrack(rtpSenderDetails.sender); err != nil {
								log.Println("error removing tracks", err)
							}

							offerConstruct := sdpConstruct{peerConnection: pc, ClassSessionID: sdp.ClassSessionID, UserID: rtpSenderDetails.userID}
							if err = offerConstruct.sendRenegotiateOffer(); err != nil {
								log.Println("failed to send renegotiation offer, closing now", err)
							}
						}
					}

					break
				}
			}

			s.audioTracks[sdp.ClassSessionID] = audioTracks

			s.peerConnectionMutexes.Unlock()
			s.audioTrackMutexes.Unlock()
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Join class session Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnSignalingStateChange(func(cc webrtc.SignalingState) {
		fmt.Println("Join class session singaling", cc.String())
	})

	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		fmt.Println("OnTrack detects join session to have", remoteTrack.PayloadType())
		audioTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sdp.ClassSessionID, sdp.UserID)
		if err != nil {
			onSessionError(user, "Unable to create an audio track.")
			peerConnection.Close()
			return
		}

		go func() {
			// Confirm both video and audio track from publisher are both enabled and publisher is still up.
			s.peerConnectionMutexes.RLock()

			publisherPeerConnection := s.peerConnection[sdp.Author]
			if publisherPeerConnection == nil {
				// Send back a JoinSessionError. Class session is closed.
				closePeerConnection(peerConnection)
				onSessionError(user, "Class session has ended.")

				s.peerConnectionMutexes.RUnlock()
				return
			}

			s.peerConnectionMutexes.RUnlock()

			s.publisherTrackMutexes.RLock()
			if s.publisherVideoTracks[sdp.ClassSessionID] == nil {
				onSessionError(user, "Publisher has not started call yet.")
				closePeerConnection(peerConnection)

				s.publisherTrackMutexes.RUnlock()
				return
			}

			// Add publishers video track.
			_, err = peerConnection.AddTrack(s.publisherVideoTracks[sdp.ClassSessionID])
			if err != nil {
				closePeerConnection(peerConnection)

				onSessionError(user, "Error adding publishers track.")
				s.publisherTrackMutexes.RUnlock()
				return
			}
			s.publisherTrackMutexes.RUnlock()

			s.peerConnectionMutexes.Lock()
			s.connectedUsersMutex.Lock()
			s.audioTrackMutexes.Lock()

			// Add other users audio tracks.
			for _, track := range s.audioTracks[sdp.ClassSessionID] {
				if track != nil {
					sender, err := peerConnection.AddTrack(track)
					if err != nil {
						log.Println("error adding audio track in join session", err)
					}

					senderData := rtpSenderData{
						userID: track.Label(), // Label has users ID.
						sender: sender}

					s.audioTrackSender[track] = append(s.audioTrackSender[track], senderData)

				} else {
					log.Println("failed track here")
				}
			}

			s.audioTracks[sdp.ClassSessionID] = append(s.audioTracks[sdp.ClassSessionID], audioTrack)

			// Send audio track to other session users and call for renegotiation.
			for _, otherConnectedUser := range s.connectedUsers[sdp.ClassSessionID] {
				pc := s.peerConnection[otherConnectedUser]

				if pc != nil && pc.ConnectionState() != webrtc.PeerConnectionStateClosed {

					sender, err := pc.AddTrack(audioTrack)
					if err != nil {
						log.Println("could not add track to other users", err)
					}

					// Save other users sender track.
					senderData := rtpSenderData{
						userID: sdp.Author,
						sender: sender}
					s.audioTrackSender[audioTrack] = append(s.audioTrackSender[audioTrack], senderData)

					offerConstruct := sdpConstruct{peerConnection: pc, ClassSessionID: sdp.ClassSessionID, UserID: otherConnectedUser}
					if err = offerConstruct.sendRenegotiateOffer(); err != nil {
						// If error is nil, there's still a chance to be corrected on the next renegotiation.
						log.Println("failed to send renegotiation offer, closing now", err)
					}
				}
			}

			// Renegotiate with self.
			offerConstruct := sdpConstruct{peerConnection: peerConnection, ClassSessionID: sdp.ClassSessionID, UserID: sdp.UserID}
			if err = offerConstruct.sendRenegotiateOffer(); err != nil {
				log.Println("failed to send renegotiation offer, closing now", err)
			}

			s.connectedUsers[sdp.ClassSessionID] = append(s.connectedUsers[sdp.ClassSessionID], sdp.UserID)
			s.peerConnection[sdp.UserID] = peerConnection

			s.audioTrackMutexes.Unlock()
			s.connectedUsersMutex.Unlock()
			s.peerConnectionMutexes.Unlock()

			fmt.Println("Starting audio writing")
		}()

		rtpBuf := make([]byte, 1400)
		for {
			i, err := remoteTrack.Read(rtpBuf)
			if err != nil {
				log.Println("error reading remote track at join session", err)
				break
			}

			_, err = audioTrack.Write(rtpBuf[:i])
			if err != nil && !errors.Is(err, io.ErrClosedPipe) {
				log.Println("Publisher video packet writed break", err)
				break
			}
		}

		log.Println("Subscriber audio track exited")
	})

	sdp.peerConnection = peerConnection
	if err = sdp.negotiate(); err != nil {
		onSessionError(user, "Unable to initiate peer negotiation.")
		log.Println("unable to negotiate on join class session", err)
	}
}

func (sdp sdpConstruct) negotiate() error {
	err := sdp.peerConnection.SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  sdp.SDP,
			Type: webrtc.SDPTypeOffer,
		})

	if err != nil {
		return err
	}

	answer, err := sdp.peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	if err := sdp.peerConnection.SetLocalDescription(answer); err != nil {
		return err
	}

	user := sdp.UserID

	sdp = sdpConstruct{}
	sdp.SDP = answer.SDP
	sdp.MsgType = "Negotiate"

	jsonContent, err := json.Marshal(sdp)
	if err != nil {
		return err
	}

	HubConstruct.sendMessage(jsonContent, user)

	return nil
}

func (sdp sdpConstruct) sendRenegotiateOffer() error {
	offer, err := sdp.peerConnection.CreateOffer(nil)
	if err != nil {
		// ToDo: Do we need to do anything further?? Can we instead broadcast to room indication fail??.
		return err
	}

	if err := sdp.peerConnection.SetLocalDescription(offer); err != nil {
		return err
	}

	data := struct {
		MsgType   string `json:"msgType"`
		SessionID string `json:"sessionID"`
		SDP       string `json:"sdp"`
	}{
		"RenegotiateSDP",
		sdp.ClassSessionID,
		offer.SDP,
	}

	jsonContent, err := json.Marshal(data)
	if err != nil {
		return err
	}

	HubConstruct.sendMessage(jsonContent, sdp.UserID)
	return nil
}

func (sdp sdpConstruct) acceptRenegotiation(msg []byte) {
	if err := json.Unmarshal(msg, &sdp); err != nil {
		log.Println("Unable to unmarshal json", err)
		return
	}

	classSessions.peerConnectionMutexes.RLock()

	peerConnection, ok := classSessions.peerConnection[sdp.UserID]
	if !ok {
		// return values.ErrPeerConnectionNotFound
		log.Println("Failed to establish accept renegotiation")
		return
	}

	classSessions.peerConnectionMutexes.RUnlock()

	err := peerConnection.SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  sdp.SDP,
			Type: webrtc.SDPTypeAnswer,
		})

	if err != nil {
		log.Println("Failed to set remote description while accepting renegotiation", err)
	}
}

func closePeerConnection(pc *webrtc.PeerConnection) {
	if pc == nil || pc.ConnectionState() == webrtc.PeerConnectionStateClosed {
		return
	}

	pc.Close()
}

func onSessionError(user, errString string) {
	errContent := struct {
		MsgType      string `json:"msgType"`
		ErrorDetails string `json:"errorDetails"`
	}{
		values.ClassSessionError,
		errString,
	}

	content, err := json.Marshal(errContent)
	if err != nil {
		log.Println("unable to send session error", err)
	}

	HubConstruct.sendMessage(content, user)
}
