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
	// Setup the codecs you want to use.
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	classSessions.api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
}

var (
	// Create a MediaEngine object to configure the supported codec
	classSessions = classSessionPeerConnections{
		publisherVideoTracks:  make(map[string]*webrtc.Track),
		publisherTrackMutexes: &sync.Mutex{},

		audioTracks:       make(map[string][]*webrtc.Track),
		audioTrackMutexes: &sync.Mutex{},

		peerConnection:        make(map[string]*webrtc.PeerConnection),
		peerConnectionMutexes: &sync.Mutex{},

		publisher:      make(map[string][]string),
		publisherMutex: &sync.Mutex{},
	}
)

func (s *classSessionPeerConnections) startClassSession(msg []byte) {
	sdp := sdpConstruct{}

	if err := json.Unmarshal(msg, &sdp); err != nil {
		return
	}

	sessionID := uuid.New().String()
	sdp.ClassSessionID = sessionID

	s.publisherMutex.Lock()
	s.publisher[sessionID] = []string{sdp.Author}
	s.publisherMutex.Unlock()

	s.peerConnectionMutexes.Lock()
	var err error
	s.peerConnection[sdp.Author], err = classSessions.api.NewPeerConnection(values.PeerConnectionConfig)
	if err != nil {
		// Send back a class session error back to user
		return
	}
	s.peerConnectionMutexes.Unlock()

	_, err = s.peerConnection[sdp.Author].AddTransceiver(webrtc.RTPCodecTypeVideo, webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	if err != nil {
		// Close connection and send back a class session error back to user
		return
	}

	_, err = s.peerConnection[sdp.Author].AddTransceiver(webrtc.RTPCodecTypeVideo, webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionSendrecv})
	if err != nil {
		// Close connection and send back a class session error back to user
		return
	}

	s.peerConnection[sdp.Author].OnConnectionStateChange(func(cc webrtc.PeerConnectionState) {
		fmt.Printf("PeerConnection State has changed %s \n", cc.String())
		if cc == webrtc.PeerConnectionStateFailed {
			// Since this is the publisher, all video and audio tracks related to the user
			// should be cleared and all peer connections closed.
			// Note: We should check if audio track is nil BEFORE creating a peerconnection
			// when joining session.
			s.peerConnectionMutexes.Lock()
			for _, user := range s.publisher[sessionID] {
				if s.peerConnection[user].ConnectionState() != webrtc.PeerConnectionStateClosed {
					s.peerConnection[user].Close()
				}

				delete(s.peerConnection, user)
			}

			s.publisherMutex.Lock()

			delete(s.publisher, sessionID)
			s.publisherMutex.Unlock()

			s.audioTrackMutexes.Lock()
			delete(s.audioTracks, sessionID)
			s.audioTrackMutexes.Unlock()

			s.publisherTrackMutexes.Lock()
			delete(s.publisherVideoTracks, sessionID)
			s.publisherTrackMutexes.Lock()

			s.peerConnectionMutexes.Unlock()
		}
	})

	s.peerConnection[sdp.Author].OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	s.peerConnection[sdp.Author].OnSignalingStateChange(func(cc webrtc.SignalingState) {
		fmt.Println("singaling", cc.String())
	})

	s.peerConnection[sdp.Author].OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Publisher <-> Server is to receive both audio and video packets.
		// Packets are to be broadcasted to other users on the session.
		// For video, resolution is in 480px.
		// ToDo: How do I confirm this???
		if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP8 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP9 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeH264 {
			videoTrack, err := s.peerConnection[sdp.Author].NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sessionID, sdp.Author)
			if err != nil {
				log.Println(err)
				// Return back a class session creation error back to client.
				// ToDo: Conver err==nil and see how peerConnectionState reacts.
				s.peerConnection[sdp.Author].Close()
				return
			}

			s.publisherTrackMutexes.Lock()
			s.publisherVideoTracks[sessionID] = videoTrack
			s.publisherTrackMutexes.Unlock()

			// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
			go func() {
				ticker := time.NewTicker(values.RtcpPLIInterval)
				for range ticker.C {
					err := s.peerConnection[sdp.Author].WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: videoTrack.SSRC()}})
					if err != nil {
						break
					}
				}
			}()

			rtpBuf := make([]byte, 1400)
			for {
				i, err := remoteTrack.Read(rtpBuf)
				if err != nil {
					log.Println(err)
					break
				}

				// ToDo: Do we really need a locker? Can we write multiple packets to a track?.
				_, err = videoTrack.Write(rtpBuf[:i])
				if err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Println("publisher video packet writed break", err)
					break
				}
			}
			log.Println("Publisher video track exited")

		} else {
			audioTrack, err := s.peerConnection[sdp.Author].NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sessionID, sdp.Author)
			if err != nil {
				log.Println(err)
				// Return back a class session creation error back to client.
				// ToDo: Convert err==nil and see how peerConnectionState reacts.
				// Also, users might decide to disable video/audio on start.
				s.peerConnection[sdp.Author].Close()
				return
			}

			s.audioTrackMutexes.Lock()
			s.audioTracks[sessionID] = []*webrtc.Track{audioTrack}
			s.audioTrackMutexes.Unlock()

			rtpBuf := make([]byte, 1400)
			for {
				i, err := remoteTrack.Read(rtpBuf)
				if err != nil {
					log.Println(err)
					break
				}

				_, err = audioTrack.Write(rtpBuf[:i])
				if err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Println("publisher video packet writed break", err)
					break
				}
			}

			log.Println("Publisher audio track exited")
		}
	})

	s.peerConnection[sdp.Author].SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  sdp.SDP,
			Type: webrtc.SDPTypeOffer,
		})

	answer, err := s.peerConnection[sdp.Author].CreateAnswer(nil)
	if err != nil {
		// ToDo: Do we need to do anything further??
		s.peerConnection[sdp.Author].Close()
		return
	}

	s.peerConnection[sdp.Author].SetLocalDescription(answer)
	sdp.SDP = answer.SDP

	jsonContent, err := json.Marshal(sdp)
	if err != nil {
		s.peerConnection[sdp.Author].Close()
		return
	}

	// Send back answer SDP to client and also class notification to all users in room.
	roomUsers, err := Message{
		RoomID:   sdp.RoomID,
		UserID:   sdp.Author,
		Type:     "classSession",
		FileHash: sdp.SDP,
	}.SaveMessageContent()

	if err != nil {
		s.peerConnection[sdp.Author].Close()
		return
	}

	for _, user := range roomUsers {
		HubConstruct.sendMessage(jsonContent, user)
	}
}

func (s *classSessionPeerConnections) joinClassSession() {

}
