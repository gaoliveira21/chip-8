package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gaoliveira21/chip8/memory"
	"github.com/gaoliveira21/chip8/utils"
)

const (
	FREQUENCY = 700 // Instructions per second

	INSTRUCTION_BITMASK = 0xF000
	X_BITMASK           = 0x0F00
	Y_BITMASK           = 0x00F0
	N_BITMASK           = 0x000F
	NN_BITMASK          = 0x00FF
	NNN_BITMASK         = 0x0FFF
)

// Display
const (
	WIDTH  = 0x40
	HEIGHT = 0x20
)

type CPU struct {
	pc         uint16   // Program Counter
	i          uint16   // I Register
	v          [16]byte // Variable registers
	mmu        memory.MMU
	display    [HEIGHT][WIDTH]byte
	delayTimer uint8
	soundTimer uint8
}

type opcode struct {
	instruction uint16
	registerX   uint8
	registerY   uint8
	n           uint8
	nn          uint8
	nnn         uint16
}

func (cpu *CPU) loadFont() {
	for i := 0x050; i <= 0x09F; i++ {
		cpu.mmu.Write(uint16(i), utils.Fontdata[i-0x050])
	}
}

func NewCpu() CPU {
	cpu := CPU{
		pc: 0x200,
	}

	cpu.loadFont()

	return cpu
}

func (cpu *CPU) LoadROM(rom string) {
	romData, err := os.ReadFile("./roms/" + rom)

	if err != nil {
		log.Fatal(err)
	}

	for index, b := range romData {
		cpu.mmu.Write(uint16(index)+0x200, b)
	}
}

func (cpu *CPU) Decode(data uint16) (oc *opcode) {
	return &opcode{
		instruction: data & INSTRUCTION_BITMASK,
		registerX:   uint8((data & X_BITMASK) >> 8),
		registerY:   uint8((data & Y_BITMASK) >> 4),
		n:           uint8(data & N_BITMASK),
		nn:          uint8(data & NN_BITMASK),
		nnn:         (data & NNN_BITMASK),
	}
}

func (cpu *CPU) Clock() {
	data := cpu.mmu.Fetch(cpu.pc)
	cpu.pc += 2

	opcode := cpu.Decode(data)

	switch opcode.instruction {
	case 0x0000:
		{
			switch opcode.n {
			case 0x0:
				cpu.cls()

			case 0xE:
				fmt.Print("RET")

			default:
				fmt.Print("sys")
			}
		}
	case 0x1000:
		cpu.jp(opcode.nnn)
	}
}

func (cpu *CPU) cls() {
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			cpu.display[i][j] = 0x00
		}
	}
}

func (cpu *CPU) jp(addr uint16) {
	cpu.pc = addr
}
