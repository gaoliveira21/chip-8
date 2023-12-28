package audio_test

import (
	"os"
	"testing"

	"github.com/gaoliveira21/chip8/core/audio"
)

func TestNewAudioPlayer(t *testing.T) {
	p, err := audio.NewAudioPlayer()

	if p == nil {
		t.Error(err)
	}
}

func TestNewAudioPlayerWeb(t *testing.T) {
	os.Setenv("EXEC_MODE", "web")

	p, _ := audio.NewAudioPlayer()

	if p != nil {
		t.Error("Audio player expected to be nil")
	}
}
