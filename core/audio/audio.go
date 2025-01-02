package audio

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed beep.mp3
var fs embed.FS

const sampleRate = 48000

type AudioPlayer = *audio.Player

func readFromFS() (io.Reader, error) {
	r, err := fs.ReadFile("beep.mp3")

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(r), nil
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
