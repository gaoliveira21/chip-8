package main

import (
	"os"
	"slices"
	"testing"

	"github.com/gaoliveira21/chip8/utils"
)

func TestNewCpu(t *testing.T) {
	cpu := NewCpu()

	inMemoryFonts := [len(utils.Fontdata)]byte{}

	for i := 0x050; i <= 0x09F; i++ {
		font := cpu.mmu.Fetch(uint16(i))

		inMemoryFonts[i-0x050] = byte(font >> 8)
	}

	if cpu.pc != 0x200 {
		t.Errorf("cpu.pc = %d; expected 0x200", cpu.pc)
	}

	if inMemoryFonts != utils.Fontdata {
		t.Error("Error loading fonts")
	}
}

func TestDecode(t *testing.T) {
	cpu := NewCpu()

	opcode := cpu.Decode(0xABCD)

	var expected uint16 = 0xA000

	if opcode.instruction != expected {
		t.Errorf("opcode.instruction = 0x%X; expected 0x%X", opcode.instruction, expected)
	}

	expected = 0xB
	if opcode.registerX != uint8(expected) {
		t.Errorf("opcode.registerX = 0x%X; expected 0x%X", opcode.registerX, expected)
	}

	expected = 0xC
	if opcode.registerY != uint8(expected) {
		t.Errorf("opcode.registerY = 0x%X; expected 0x%X", opcode.registerY, expected)
	}

	expected = 0xD
	if opcode.n != uint8(expected) {
		t.Errorf("opcode.n = 0x%X; expected 0x%X", opcode.n, expected)
	}

	expected = 0xCD
	if opcode.nn != uint8(expected) {
		t.Errorf("opcode.n = 0x%X; expected 0x%X", opcode.nn, expected)
	}

	expected = 0xBCD
	if opcode.nnn != expected {
		t.Errorf("opcode.n = 0x%X; expected 0x%X", opcode.nnn, expected)
	}
}

func TestLoadROM(t *testing.T) {
	originalROMData, err := os.ReadFile("./roms/IBM.ch8")

	if err != nil {
		t.Fatal(err)
	}

	cpu := NewCpu()

	cpu.LoadROM("IBM.ch8")

	inMemoryROM := []byte{}

	for i := 0; i < len(originalROMData); i++ {
		romByte := cpu.mmu.Fetch(uint16(i + 0x200))

		inMemoryROM = append(inMemoryROM, byte(romByte>>8))
	}

	if !slices.Equal[[]byte](inMemoryROM, originalROMData) {
		t.Error("Error loading ROM")
	}
}

func TestClock(t *testing.T) {
	cpu := NewCpu()

	currentPc := cpu.pc

	cpu.Clock()

	expected := currentPc + 2
	if cpu.pc != expected {
		t.Errorf("cpu.pc = 0x%X; expected 0x%X", cpu.pc, expected)
	}
}

func TestCLS(t *testing.T) {
	cpu := NewCpu()

	cpu.display[0][0] = 0xFF
	cpu.display[0][1] = 0xEF

	cpu.cls()

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			if cpu.display[i][j] != 0x00 {
				t.Errorf("cpu.display[%d][%d] = 0x%X; expected 0x00", i, j, cpu.display[i][j])
			}
		}
	}
}

func TestJP(t *testing.T) {
	cpu := NewCpu()

	expected := 0xFFF

	cpu.jp(uint16(expected))

	if cpu.pc != uint16(expected) {
		t.Errorf("cpu.pc = 0x%X; expected 0x%X", cpu.pc, expected)
	}
}

func TestLD(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1
	var expected uint8 = 0xFF

	cpu.ld(vIndex, expected)

	if cpu.v[vIndex] != expected {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}
}

func TestADD(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1

	cpu.add(vIndex, 0x02)
	cpu.add(vIndex, 0x03)

	expected := 0x05

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}
}

func TestLDI(t *testing.T) {
	cpu := NewCpu()

	var expected uint16 = 0x0ABC

	cpu.ldi(expected)

	if cpu.i != expected {
		t.Errorf("cpu.i = 0x%X; expected 0x%X", cpu.i, expected)
	}
}
