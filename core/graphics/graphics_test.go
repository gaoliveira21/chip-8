package graphics_test

import (
	"testing"

	"github.com/gaoliveira21/chip8/core/graphics"
)

func TestSetAndGetPixel(t *testing.T) {
	g := graphics.NewGraphics()

	g.SetPixel(0, 0, 0x1)

	if g.GetPixel(0, 0) != 0x1 {
		t.Errorf("graphics.Display[0][0] = 0x%X; expected 0x01", g.GetPixel(0, 0))
	}
}
