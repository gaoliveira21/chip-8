package main

import (
	"github.com/gaoliveira21/chip8/utils"
)

const RAM_SIZE = 4096 // 4 KB
const FREQUENCY = 700 // Instructions per second

// Display
const (
	WIDTH  = 0x40
	HEIGHT = 0x20
)

type CPU struct {
	memory     [RAM_SIZE]uint8
	pc         uint16 // Program Counter
	i          uint16 // I Register
	v          uint8  // Variable registers
	stack      utils.Stack
	display    [HEIGHT][WIDTH]byte
	delayTimer uint8
	soundTimer uint8
}

func (cpu *CPU) LoadFont() {
	fontIndex := 0
	for i := 0x050; i <= 0x09F; i++ {
		cpu.memory[i] = fontdata[fontIndex]
		fontIndex++
	}
}
