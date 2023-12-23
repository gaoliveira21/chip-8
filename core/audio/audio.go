package audio

import (
	"errors"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 48000

type AudioPlayer = *audio.Player

func readFromFS() (io.Reader, error) {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "..", "..", "cli", "assets", "beep.mp3")

	r, err := os.Open(d)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewAudioPlayer() (AudioPlayer, error) {
	execMode := os.Getenv("EXEC_MODE")

	var r io.Reader
	var err error

	if execMode == "web" {
		return nil, errors.New("audio unsuported on web")
	}

	r, err = readFromFS()

	if err != nil {
		return nil, err
	}

	stream, err := mp3.DecodeWithSampleRate(sampleRate, r)

	if err != nil {
		return nil, err
	}

	audioContext := audio.NewContext(sampleRate)
	p, err := audioContext.NewPlayer(stream)

	if err != nil {
		return nil, err
	}

	return p, nil
}
