package audio

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 48000

type AudioPlayer = *audio.Player

func NewAudioPlayer(mp3FilePath string) (AudioPlayer, error) {
	audioContext := audio.NewContext(sampleRate)
	r, err := os.Open(mp3FilePath)

	if err != nil {
		return nil, err
	}

	stream, err := mp3.DecodeWithSampleRate(sampleRate, r)

	if err != nil {
		return nil, err
	}

	p, err := audioContext.NewPlayer(stream)

	if err != nil {
		return nil, err
	}

	return p, nil
}
