package main

import (
	"io"
	"log"
	"time"

	"github.com/ebml-go/webm"
)

type Stream interface {
	Meta() *webm.WebM
	VideoDecoder() *VideoDecoder
	AudioDecoder() *AudioDecoder
	Seek(d time.Duration)
}

type webmStream struct {
	meta webm.WebM
	vdec *VideoDecoder
	adec *AudioDecoder

	reader *webm.Reader
}

func NewStream(r io.ReadSeeker) (Stream, error) {
	s := &webmStream{}
	reader, err := webm.Parse(r, &s.meta)
	if err != nil {
		return nil, err
	}
	s.reader = reader
	vtrack := s.meta.FindFirstVideoTrack()
	atrack := s.meta.FindFirstAudioTrack()
	vPackets := make(chan webm.Packet, 32)
	aPackets := make(chan webm.Packet, 32)
	if vtrack != nil {
		log.Printf("webm: found video track: %dx%d dur: %v %s", vtrack.DisplayWidth,
			vtrack.DisplayHeight, s.meta.Segment.GetDuration(), vtrack.CodecID)

		s.vdec, err = NewVideoDecoer(VCodec(vtrack.CodecID), vPackets)
		if err != nil {
			return nil, err
		}
	}
	if atrack != nil {
		log.Printf("webm: found audio track: ch: %d %.1fHz, dur: %v, codec: %s", atrack.Channels,
			atrack.SamplingFrequency, s.meta.Segment.GetDuration(), atrack.CodecID)

		s.adec, err = NewAudioDecoder(AudioCodec(atrack.CodecID), atrack.CodecPrivate,
			int(atrack.Channels), int(atrack.SamplingFrequency), aPackets)
		if err != nil {
			return nil, err
		}
	}
	go func() { // demuxer
		for pkt := range s.reader.Chan {
			switch {
			case vtrack == nil:
				aPackets <- pkt // audio only
			case atrack == nil:
				vPackets <- pkt // video only
			default:
				switch pkt.TrackNumber {
				case vtrack.TrackNumber:
					vPackets <- pkt
				case atrack.TrackNumber:
					aPackets <- pkt
				}
			}
		}
		close(vPackets)
		close(aPackets)
		s.reader.Shutdown()
	}()
	return s, nil
}

func (s *webmStream) Meta() *webm.WebM {
	return &s.meta
}

func (s *webmStream) VideoDecoder() *VideoDecoder {
	return s.vdec
}

func (s *webmStream) AudioDecoder() *AudioDecoder {
	return s.adec
}

func (s *webmStream) Seek(d time.Duration) {
	s.reader.Seek(d)
}
