package main

import (
	"fmt"
	"image"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ebml-go/webm"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xlab/libvpx-go/vpx"
)

type Frame struct {
	*image.RGBA
	Timecode   time.Duration
	IsKeyframe bool
}

type VideoDecoder struct {
	src   <-chan webm.Packet
	ctx   *vpx.CodecCtx
	iface *vpx.CodecIface

	offscreen *ebiten.Image
	m         sync.Mutex

	pos atomic.Int64
}

type VCodec string

const (
	CodecVP8  VCodec = "V_VP8"
	CodecVP9  VCodec = "V_VP9"
	CodecVP10 VCodec = "V_VP10"
)

type Positioner interface {
	Position() time.Duration
}

func NewVideoDecoer(codec VCodec, src <-chan webm.Packet) (*VideoDecoder, error) {
	dec := &VideoDecoder{
		src: src,
		ctx: vpx.NewCodecCtx(),
	}
	switch codec {
	case CodecVP8:
		dec.iface = vpx.DecoderIfaceVP8()
	case CodecVP9:
		dec.iface = vpx.DecoderIfaceVP9()
	default:
		return nil, fmt.Errorf("webm-player: unsupported VPX codec: %s", codec)
	}
	if err := vpx.Error(vpx.CodecDecInitVer(dec.ctx, dec.iface, nil, 0, vpx.DecoderABIVersion)); err != nil {
		return nil, err
	}
	go dec.loop()
	return dec, nil
}

type VideoDecoderDrawOptions struct {
	GeoM       ebiten.GeoM
	ColorScale ebiten.ColorScale
	Blend      ebiten.Blend
}

func (v *VideoDecoder) Update(position time.Duration) {
	v.pos.Store(int64(position))
}

func (v *VideoDecoder) Draw(dst *ebiten.Image, options *VideoDecoderDrawOptions) {
	v.m.Lock()
	defer v.m.Unlock()

	if v.offscreen == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	if options != nil {
		op.GeoM = options.GeoM
		op.ColorScale = options.ColorScale
		op.Blend = options.Blend
	}
	dst.DrawImage(v.offscreen, op)
}

func (v *VideoDecoder) loop() {
loop:
	for pkt := range v.src {
		dataSize := uint32(len(pkt.Data))
		if err := vpx.Error(vpx.CodecDecode(v.ctx, string(pkt.Data), dataSize, nil, 0)); err != nil {
			log.Println("[WARN]", err)
			continue
		}
		pos := time.Duration(v.pos.Load())
		if pos-time.Second/60 > pkt.Timecode {
			continue loop
		}

		var iter vpx.CodecIter
		for img := vpx.CodecGetFrame(v.ctx, &iter); img != nil; img = vpx.CodecGetFrame(v.ctx, &iter) {
			img.Deref()
			if pos < pkt.Timecode {
				time.Sleep(pkt.Timecode - pos)
			}
			// TODO: Use img.ImageYCbCr and a shader.
			img := img.ImageRGBA()

			v.m.Lock()
			if v.offscreen != nil && v.offscreen.Bounds() != img.Bounds() {
				v.offscreen.Deallocate()
				v.offscreen = nil
			}
			if v.offscreen == nil {
				v.offscreen = ebiten.NewImage(img.Bounds().Dx(), img.Bounds().Dy())
			}
			v.offscreen.WritePixels(img.Pix)
			v.m.Unlock()
		}
	}
}
