package core

import (
	"image/color"
	"log"

	"github.com/gaoliveira21/chip8/core/audio"
	"github.com/gaoliveira21/chip8/core/cpu"
	"github.com/gaoliveira21/chip8/core/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chip8 struct {
	cpu         *cpu.CPU
	square      *ebiten.Image
	audioPlayer audio.AudioPlayer
}

func (c8 *Chip8) Update() error {
	for i := 0; i < cpu.SPEED; i++ {
		for key, value := range input.Keypad {
			if ebiten.IsKeyPressed(key) {
				c8.cpu.Keys[value] = 0x01
			} else {
				c8.cpu.Keys[value] = 0x00
			}
		}

		c8.cpu.Run()

		if c8.cpu.ResizeWindow {
			ebiten.SetWindowSize(c8.cpu.Graphics.Width*10, c8.cpu.Graphics.Height*10)
			c8.cpu.ResizeWindow = false
		}

		if c8.cpu.SoundTimer > 0 && c8.audioPlayer != nil {
			c8.audioPlayer.Play()
			c8.audioPlayer.Rewind()
		}
	}

	return nil
}

func (c8 *Chip8) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{23, 20, 33, 1})

	for h := 0; h < c8.cpu.Graphics.Height; h++ {
		for w := 0; w < c8.cpu.Graphics.Width; w++ {
			if c8.cpu.Graphics.GetPixel(h, w) == 0x01 {
				imgOpts := &ebiten.DrawImageOptions{}
				imgOpts.GeoM.Translate(float64(w*10), float64(h*10))
				screen.DrawImage(c8.square, imgOpts)
			}
		}
	}
}

func (c8 *Chip8) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return c8.cpu.Graphics.Width * 10, c8.cpu.Graphics.Height * 10
}

func RunChip8(rom []byte, title string) {
	c := cpu.NewCpu()
	c.LoadROM(rom)

	sqr := ebiten.NewImage(10, 10)
	sqr.Fill(color.RGBA{51, 209, 122, 1})

	p, err := audio.NewAudioPlayer()

	if err != nil {
		log.Print(err)
	}

	c8 := &Chip8{
		square:      sqr,
		cpu:         &c,
		audioPlayer: p,
	}

	ebiten.SetWindowSize(c.Graphics.Width*10, c.Graphics.Height*10)
	ebiten.SetWindowTitle(title)

	if err := ebiten.RunGame(c8); err != nil {
		log.Fatal(err)
	}
}
