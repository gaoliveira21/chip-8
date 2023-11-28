package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Chip8 struct {
	cpu          *CPU
	square       *ebiten.Image
	screenWidth  int
	screenHeight int
}

func (c8 *Chip8) Update() error {
	c8.cpu.Run()

	return nil
}

func (c8 *Chip8) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})

	for h := 0; h < HEIGHT; h++ {
		for w := 0; w < WIDTH; w++ {
			if c8.cpu.display[h][w] == 0x01 {
				imgOpts := &ebiten.DrawImageOptions{}
				imgOpts.GeoM.Translate(float64(w*10), float64(h*10))
				screen.DrawImage(c8.square, imgOpts)
			}
		}
	}
}

func (c8 *Chip8) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return c8.screenWidth, c8.screenHeight
}

func main() {
	cpu := NewCpu()
	cpu.LoadROM("IBM.ch8")

	sqr := ebiten.NewImage(10, 10)
	sqr.Fill(color.White)

	c8 := &Chip8{
		square:       sqr,
		cpu:          &cpu,
		screenWidth:  WIDTH * 10,
		screenHeight: HEIGHT * 10,
	}

	ebiten.SetWindowSize(c8.screenWidth, c8.screenHeight)
	ebiten.SetWindowTitle("IBM.ch8")

	if err := ebiten.RunGame(c8); err != nil {
		log.Fatal(err)
	}
}
