package audio_test

import (
	"testing"

	"github.com/gaoliveira21/chip8/core/audio"
)

func Test(t *testing.T) {
	p, err := audio.NewAudioPlayer()

	if p == nil {
		t.Error(err)
	}
}
