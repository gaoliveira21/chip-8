package main

import (
	"image/color"
	"log"

	"github.com/gaoliveira21/chip8/core"
	"github.com/gaoliveira21/chip8/core/audio"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chip8 struct {
	cpu          *core.CPU
	square       *ebiten.Image
	audioPlayer  audio.AudioPlayer
	screenWidth  int
	screenHeight int
}

func (c8 *Chip8) Update() error {
	for i := 0; i < int(core.FREQUENCY/60); i++ {
		for key, value := range core.Keypad {
			if ebiten.IsKeyPressed(key) {
				c8.cpu.Keys[value] = 0x01
			} else {
				c8.cpu.Keys[value] = 0x00
			}
		}

		c8.cpu.Run()

		if c8.cpu.SoundTimer > 0 {
			c8.audioPlayer.Play()
			c8.audioPlayer.Rewind()
		}
	}

	return nil
}

func (c8 *Chip8) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})

	for h := 0; h < core.HEIGHT; h++ {
		for w := 0; w < core.WIDTH; w++ {
			if c8.cpu.Display[h][w] == 0x01 {
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
	cpu := core.NewCpu()
	cpu.LoadROM("./roms/PONG.ch8")

	sqr := ebiten.NewImage(10, 10)
	sqr.Fill(color.White)

	p, err := audio.NewAudioPlayer("assets/beep.mp3")

	if err != nil {
		log.Fatal(err)
	}

	c8 := &Chip8{
		square:       sqr,
		cpu:          &cpu,
		audioPlayer:  p,
		screenWidth:  core.WIDTH * 10,
		screenHeight: core.HEIGHT * 10,
	}

	ebiten.SetWindowSize(c8.screenWidth, c8.screenHeight)
	ebiten.SetWindowTitle("PONG.ch8")

	if err := ebiten.RunGame(c8); err != nil {
		log.Fatal(err)
	}
}
