package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xlab/closer"
)

var maxFps = flag.Int("fps", 30, "Limits the rendering FPS rate. Set this to 60fps for 720p60 videos")

const appName = "WebM VP8/VP9 Player"

var rateLimitDur time.Duration

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "A simple WebM player with support of VP8/VP9 video and Vorbis/Opus audio. Version: v1.0rc1")
		fmt.Fprintf(os.Stderr, "Usage: %s <file1.webm> [file2.webm]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Specify files to read streams from, sometimes audio is stored in a separate file, use the optional argument for that.")
		flag.PrintDefaults()
	}
	flag.Parse()
	rateLimitDur = time.Second / time.Duration(*maxFps)
	runtime.LockOSThread()
}

func main() {
	// defer closer.Close()
	closer.Bind(func() {
		log.Println("Bye!")
	})

	// Open WebM files
	streams := make([]io.ReadSeeker, 0, 2)
	for _, opt := range flag.Args() {
		f, err := os.Open(opt)
		if err != nil {
			log.Println("[ERR] failed to open file:", err)
		}
		streams = append(streams, f)
		if len(streams) >= 2 {
			break
		}
	}
	stream1, stream2 := discoverStreams(streams...)
	if stream1 == nil {
		closer.Fatalln("[ERR] nothing to play")
	}

	aDec := stream1.AudioDecoder()
	if stream2 != nil {
		aDec = stream2.AudioDecoder()
	}

	var view *View
	if vtrack := stream1.Meta().FindFirstVideoTrack(); vtrack != nil {
		v, err := NewView(vtrack.DisplayWidth, vtrack.DisplayHeight, stream1.VideoDecoder(), aDec)
		if err != nil {
			log.Println("[ERR] failed to create view:", err)
			return
		}
		view = v
	} else {
		v, err := NewView(0, 0, nil, aDec)
		if err != nil {
			log.Println("[ERR] failed to create view:", err)
			return
		}
		view = v
	}
	view.SetOnSeek(func(d time.Duration) {
		stream1.Seek(d)
		if stream2 != nil {
			stream2.Seek(d)
		}
	})

	ebiten.SetWindowSize(winWidth, winHeight)
	if err := ebiten.RunGame(view); err != nil {
		log.Println("[ERR] GUI loop:", err)
	}
}

// discoverStreams returns both Video and Audio streams if in separate inputs,
// otherwise only the first stream would be returned (V/A/V+A).
func discoverStreams(streams ...io.ReadSeeker) (Stream, Stream) {
	if len(streams) == 0 {
		log.Println("[WARN] no streams found")
		return nil, nil
	}

	if len(streams) == 1 {
		stream, err := NewStream(streams[0])
		if err != nil {
			log.Println("[WARN] failed to open stream:", err)
			return nil, nil
		}
		return stream, nil
	}

	var stream1Video bool
	var stream1Audio bool
	stream1, err := NewStream(streams[0])
	if err == nil {
		stream1Video = stream1.Meta().FindFirstVideoTrack() != nil
		stream1Audio = stream1.Meta().FindFirstAudioTrack() != nil
	} else {
		log.Println("[WARN] failed to open the first stream:", err)
	}
	if stream1Video && stream1Audio {
		log.Println("[INFO] found both Video+Audio in the first stream")
		return stream1, nil
	}
	var stream2Video bool
	var stream2Audio bool
	stream2, err := NewStream(streams[1])
	if err == nil {
		stream2Video = stream2.Meta().FindFirstVideoTrack() != nil
		stream2Audio = stream2.Meta().FindFirstAudioTrack() != nil
	} else {
		log.Println("[WARN] failed to open the second stream:", err)
	}
	switch {
	case stream1Video && stream2Audio:
		log.Println("[INFO] took Video from the first stream, Audio from the second")
		return stream1, stream2
	case stream1Audio && stream2Video:
		log.Println("[INFO] took Audio from the first stream, Video from the second")
		return stream2, stream1
	case stream1Video:
		log.Println("[INFO] took Video from the first stream, no Audio found")
		return stream1, nil
	case stream2Video:
		log.Println("[INFO] took Video from the second stream, no Audio found")
		return stream2, nil
	case stream1Audio:
		log.Println("[INFO] took Audio from the first stream, no Video found")
		return stream1, nil
	case stream2Audio:
		log.Println("[INFO] took Audio from the second stream, no Video found")
		return stream2, nil
	default:
		log.Println("[INFO] neither of Video or Audio found")
		return nil, nil
	}
}
