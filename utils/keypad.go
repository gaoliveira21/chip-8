package utils

import "github.com/hajimehoshi/ebiten/v2"

var Keypad = map[ebiten.Key]uint8{
	ebiten.Key1: 0x01,
	ebiten.Key2: 0x02,
	ebiten.Key3: 0x03,
	ebiten.Key4: 0x0C,
	ebiten.KeyQ: 0x04,
	ebiten.KeyW: 0x05,
	ebiten.KeyE: 0x06,
	ebiten.KeyR: 0x0D,
	ebiten.KeyA: 0x07,
	ebiten.KeyS: 0x08,
	ebiten.KeyD: 0x09,
	ebiten.KeyF: 0x0E,
	ebiten.KeyZ: 0x0A,
	ebiten.KeyX: 0x00,
	ebiten.KeyC: 0x0B,
	ebiten.KeyV: 0x0F,
}
