package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/ebml-go/webm"
	"github.com/xlab/opus-go/opus"
	"github.com/xlab/vorbis-go/decoder"
	"github.com/xlab/vorbis-go/vorbis"
)

const samplesPerBuffer = 1024

type AudioDecoder struct {
	codec      AudioCodec
	channels   int
	sampleRate int

	src     <-chan webm.Packet
	packets []webm.Packet

	voDSP   vorbis.DspState
	voBlock vorbis.Block
	voPCM   [][][]float32

	opDecoder *opus.Decoder
	opPCM     []float32

	frames []float32
}

type AudioCodec string

const (
	CodecVorbis AudioCodec = "A_VORBIS"
	CodecOpus   AudioCodec = "A_OPUS"
)

func NewAudioDecoder(codec AudioCodec, codecPrivate []byte, channels, sampleRate int, src <-chan webm.Packet) (*AudioDecoder, error) {
	d := &AudioDecoder{
		channels:   channels,
		sampleRate: sampleRate,
		codec:      codec,
		src:        src,
	}
	switch codec {
	case CodecVorbis:
		var info vorbis.Info
		vorbis.InfoInit(&info)
		var comment vorbis.Comment
		vorbis.CommentInit(&comment)
		err := decoder.ReadHeaders(codecPrivate, &info, &comment)
		if err != nil {
			return nil, err
		}
		info.Deref()
		comment.Deref()
		if comment.Comments > 0 {
			comment.UserComments = make([][]byte, comment.Comments)
			comment.Deref()
			streamInfo := decoder.ReadInfo(&info, &comment)
			log.Println("vorbis:", streamInfo.Comments)
		}
		if int(info.Channels) != channels {
			d.channels = int(channels)
			log.Printf("[WARN] vorbis: channel count mismatch %d != %d", info.Channels, channels)
		}
		if int(info.Rate) != sampleRate {
			d.sampleRate = int(info.Rate)
			log.Printf("[WARN] vorbis: sample rate mismatch %d != %d", info.Rate, sampleRate)
		}
		ret := vorbis.SynthesisInit(&d.voDSP, &info)
		if ret != 0 {
			return nil, fmt.Errorf("webm-player: vorbis.SynthesisInit failed: %d", ret)
		}
		d.voPCM = [][][]float32{
			make([][]float32, channels),
		}
		vorbis.BlockInit(&d.voDSP, &d.voBlock)
		return d, nil
	case CodecOpus:
		var err int32
		d.opDecoder = opus.DecoderCreate(int32(sampleRate), int32(channels), &err)
		if err != opus.Ok {
			return nil, fmt.Errorf("webm-player: opus.DecoderCreate failed: %d", err)
		}
		d.opPCM = make([]float32, samplesPerBuffer*channels)
		return d, nil
	default:
		return d, fmt.Errorf("webm-player: unsupported audio codec: %s", codec)
	}
}

func (a *AudioDecoder) Read(buf []byte) (int, error) {
readFrames:
	if len(a.frames) > 0 {
		n := copy(unsafe.Slice((*float32)(unsafe.Pointer(unsafe.SliceData(buf))), len(buf)/4), a.frames)
		a.frames = a.frames[n:]
		return 4 * n, nil
	}

	for len(a.packets) == 0 {
		pkt, ok := <-a.src
		if !ok {
			n := min(len(buf)/4*4, 256)
			for i := range n {
				buf[i] = 0
			}
			return n, nil
		}
		if len(pkt.Data) == 0 {
			continue
		}
		a.packets = append(a.packets, pkt)
	}

	pkt := a.packets[0]
	a.packets = a.packets[1:]

	switch a.codec {
	case CodecVorbis:
		packet := &vorbis.OggPacket{
			Packet: pkt.Data,
			Bytes:  len(pkt.Data),
		}
		if ret := vorbis.Synthesis(&a.voBlock, packet); ret != 0 {
			return 0, fmt.Errorf("webm-player: vorbis.Synthesis failed: %d", ret)
		}

		vorbis.SynthesisBlockin(&a.voDSP, &a.voBlock)

		sampleCount := vorbis.SynthesisPcmout(&a.voDSP, a.voPCM)
		if sampleCount == 0 {
			vorbis.SynthesisRead(&a.voDSP, sampleCount)
			return 0, nil
		}

		for ; sampleCount > 0; sampleCount = vorbis.SynthesisPcmout(&a.voDSP, a.voPCM) {
			for i := 0; i < int(sampleCount); i++ {
				for j := 0; j < a.channels; j++ {
					v := a.voPCM[0][j][:sampleCount][i]
					a.frames = append(a.frames, v)
					if a.channels == 1 {
						a.frames = append(a.frames, v)
					}
				}
			}
			vorbis.SynthesisRead(&a.voDSP, sampleCount)
		}

		goto readFrames

	case CodecOpus:
		sampleCount := opus.DecodeFloat(a.opDecoder, pkt.Data, int32(len(pkt.Data)), a.opPCM, samplesPerBuffer, 0)
		if sampleCount <= 0 {
			return 0, nil
		}

		origLen := len(a.frames)
		a.frames = append(a.frames, a.opPCM[:int(sampleCount)*a.channels]...)
		if a.channels == 1 {
			a.frames = append(a.frames, make([]float32, sampleCount)...)
			frames := a.frames[origLen:]
			for i := int(sampleCount) - 1; i > 0; i-- {
				frames[2*i] = frames[i]
				frames[2*i+1] = frames[i]
			}
		}

		goto readFrames

	default:
		return 0, fmt.Errorf("webm-player: unsupported audio codec: %s", a.codec)
	}
}

func (a *AudioDecoder) Channels() int {
	return a.channels
}

func (a *AudioDecoder) SampleRate() int {
	return a.sampleRate
}
