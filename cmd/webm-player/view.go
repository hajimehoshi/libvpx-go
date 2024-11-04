package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	winWidth  = 800
	winHeight = 500
)

type View struct {
	bgImg *ebiten.Image

	width  uint
	height uint

	onPause func()
	onSeek  func(d time.Duration)

	videoDecoder *VideoDecoder
	audioDecoder *AudioDecoder
	audioPlayer  *audio.Player
}

func NewView(width, height uint, videoDecoder *VideoDecoder, audioDecoder *AudioDecoder) (*View, error) {
	v := &View{
		width:        width,
		height:       height,
		videoDecoder: videoDecoder,
		audioDecoder: audioDecoder,
	}
	if audioDecoder != nil {
		ctx := audio.NewContext(audioDecoder.SampleRate())
		p, err := ctx.NewPlayerF32(audioDecoder)
		if err != nil {
			return nil, err
		}
		p.Play()
		v.audioPlayer = p
	}

	return v, nil
}

func (v *View) Update() error {
	v.videoDecoder.Update(v.audioPlayer.Position())
	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	if v.bgImg != nil {
		screen.DrawImage(v.bgImg, nil)
	}
	if v.videoDecoder != nil {
		op := &VideoDecoderDrawOptions{}
		scale := min(float64(winWidth)/float64(v.width), float64(winHeight)/float64(v.height))
		op.GeoM.Scale(scale, scale)
		v.videoDecoder.Draw(screen, op)
	}
}

func (v *View) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}

func (v *View) SetOnPause(fn func()) {
	v.onPause = fn
}

func (v *View) SetOnSeek(fn func(d time.Duration)) {
	v.onSeek = fn
}
