package model

import (
	"fmt"
	"log"
	"os"

	"github.com/at-wat/ebml-go/webm"
	"github.com/metaclips/LetsTalk/values"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v2/pkg/media/samplebuilder"
)

// newWebmWriter writes video class session to either be upload to a Dropbox drive
// or if a token is not specified, saved to mongodb using gridFS.
// ToDo: Allow users access to download from DB if file is saved using gridFS.
func newWebmWriter(fileName string) *webmWriter {
	return &webmWriter{
		fileName:     fileName,
		audioBuilder: samplebuilder.New(10, &codecs.OpusPacket{}),
		videoBuilder: samplebuilder.New(10, &codecs.VP8Packet{}),
	}
}

func (s *webmWriter) initWriter(width, height int) {
	w, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println(err)
	}

	ws, err := webm.NewSimpleBlockWriter(w,
		[]webm.TrackEntry{
			{
				Name:            "Audio",
				TrackNumber:     1,
				TrackUID:        12345,
				CodecID:         "A_OPUS",
				TrackType:       2,
				DefaultDuration: 20000000,
				Audio: &webm.Audio{
					SamplingFrequency: 48000.0,
					Channels:          2,
				},
			}, {
				Name:            "Video",
				TrackNumber:     2,
				TrackUID:        67890,
				CodecID:         "V_VP8",
				TrackType:       1,
				DefaultDuration: 33333333,
				Video: &webm.Video{
					PixelWidth:  uint64(width),
					PixelHeight: uint64(height),
				},
			},
		})
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("WebM saver has started with video width=%d, height=%d\n", width, height)
	s.audioWriter = ws[0]
	s.videoWriter = ws[1]
}

// Append Audi file.
func (s *webmWriter) pushOpus(rtpPacket *rtp.Packet) {
	s.audioBuilder.Push(rtpPacket)

	for {
		sample := s.audioBuilder.Pop()
		if sample == nil {
			return
		}

		if s.audioWriter != nil {
			s.audioTimestamp += sample.Samples
			t := s.audioTimestamp / 48
			if _, err := s.audioWriter.Write(true, int64(t), sample.Data); err != nil {
				log.Println(err)
			}
		}
	}
}

// Push to video.
func (s *webmWriter) pushVP8(rtpPacket *rtp.Packet) {
	s.videoBuilder.Push(rtpPacket)

	for {
		sample := s.videoBuilder.Pop()
		if sample == nil {
			return
		}

		// Read VP8 header.
		videoKeyframe := (sample.Data[0]&0x1 == 0)
		if videoKeyframe {
			// Keyframe has frame information.
			raw := uint(sample.Data[6]) | uint(sample.Data[7])<<8 | uint(sample.Data[8])<<16 | uint(sample.Data[9])<<24
			width := int(raw & 0x3FFF)
			height := int((raw >> 16) & 0x3FFF)

			if s.videoWriter == nil || s.audioWriter == nil {
				// Initialize WebM saver using received frame size.
				s.initWriter(width, height)
			}
		}

		if s.videoWriter != nil {
			s.videoTimestamp += sample.Samples
			t := s.videoTimestamp / 90
			if _, err := s.videoWriter.Write(videoKeyframe, int64(t), sample.Data); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *webmWriter) close() {
	fmt.Printf("Finalizing webm...\n")
	if s.audioWriter != nil {
		if err := s.audioWriter.Close(); err != nil {
			log.Println(err)
		}
	}

	if s.videoWriter != nil {
		if err := s.videoWriter.Close(); err != nil {
			log.Println(err)
		}
	}

	log.Println("Video writer closed for session, uploading.", s.fileName)
}

func (s *webmWriter) getVideoFileSharableLink() (string, error) {
	if values.Config.DropboxToken != "" {
		dropBoxUploader, err := newDropboxUploader(s.fileName)
		if err != nil {
			log.Println("unable to initialize dropbox uploader", err)
			return "", err
		}

		link, err := dropBoxUploader.dropboxFileUploader()
		if err != nil {
			log.Println("unable to get dropbox sharable link", err)
			return "", err
		}

		log.Println("file uploaded to file server")

		return link, nil
	}

	return "", nil
}

func (s *webmWriter) uploadToDB() {
	defer func() {
		if err := os.Remove(s.fileName); err != nil {
			log.Println("unable to remove file", err)
		}
	}()

	if err := uploadFileGridFS(s.fileName); err != nil {
		log.Println("error saving file to DB", err)
		return
	}

	log.Println("File uploaded to DB")
}
