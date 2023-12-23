package audio

import (
	"bytes"
	"io"
	"os"

	"github.com/gaoliveira21/chip8/core/http"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 48000

type AudioPlayer = *audio.Player

func readFromServer(mp3FilePath string) (io.Reader, error) {
	b, err := http.ReadFile(mp3FilePath)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func readFromFS(mp3FilePath string) (io.Reader, error) {
	r, err := os.Open(mp3FilePath)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewAudioPlayer(mp3FilePath string) (AudioPlayer, error) {
	execMode := os.Getenv("EXEC_MODE")

	var r io.Reader
	var err error

	if execMode == "web" {
		r, err = readFromServer(mp3FilePath)
	} else {
		r, err = readFromFS(mp3FilePath)
	}

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
