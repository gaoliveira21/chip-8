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
	// TODO: Improve clock test with real memory addresses
	t.Skip()
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

func TestRET(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Stack.Push(0xDDEE)
	cpu.ret()

	if cpu.pc != 0xDDEE {
		t.Errorf("cpu.pc = 0x%X; expected 0xDDEE", cpu.pc)
	}
}

func TestCALL(t *testing.T) {
	cpu := NewCpu()

	cpu.pc = 0x300
	cpu.call(0x400)

	stackPC := cpu.mmu.Stack.Pop()
	currentPC := cpu.pc

	if stackPC != 0x300 {
		t.Errorf("Stack PC = 0x%X; expected 0x300", stackPC)
	}

	if currentPC != 0x400 {
		t.Errorf("Current PC = 0x%X; expected 0x400", currentPC)
	}
}

func TestJPWithoutOffset(t *testing.T) {
	cpu := NewCpu()

	expected := 0xFFF

	cpu.jp(uint16(expected), 0)

	if cpu.pc != uint16(expected) {
		t.Errorf("cpu.pc = 0x%X; expected 0x%X", cpu.pc, expected)
	}
}

func TestJPWithOffset(t *testing.T) {
	cpu := NewCpu()
	cpu.v[0x0] = 0x02

	var addr uint16 = 0xFF0
	expected := addr + uint16(cpu.v[0x0])

	cpu.jp(uint16(addr), cpu.v[0x0])

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

func TestLDT(t *testing.T) {
	cpu := NewCpu()

	cpu.ldt(0x60)

	if cpu.delayTimer != 0x60 {
		t.Errorf("cpu.delayTimer = 0x%X; expected 0x60", cpu.delayTimer)
	}
}

func TestLDS(t *testing.T) {
	cpu := NewCpu()

	cpu.lds(0x80)

	if cpu.soundTimer != 0x80 {
		t.Errorf("cpu.soundTimer = 0x%X; expected 0x80", cpu.soundTimer)
	}
}

func TestLDKWithNoKeyPressed(t *testing.T) {
	cpu := NewCpu()
	cpu.pc += 2

	var vIndex uint8 = 0x1

	cpu.ldk(vIndex)

	if cpu.v[vIndex] != 0x0 {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], 0x0)
	}

	if cpu.pc != 0x200 {
		t.Errorf("cpu.pc = 0x%X; expected 0x200", cpu.pc)
	}
}

func TestLDKWithKeyPressed(t *testing.T) {
	cpu := NewCpu()
	cpu.pc += 2

	var vIndex uint8 = 0x1
	var keyPressed uint8 = 0xF

	cpu.Keys[keyPressed] = 0x1

	cpu.ldk(vIndex)

	if cpu.v[vIndex] != keyPressed {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], keyPressed)
	}

	if cpu.pc != 0x202 {
		t.Errorf("cpu.pc = 0x%X; expected 0x200", cpu.pc)
	}
}

func TestADDWithoutCarry(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1

	cpu.add(vIndex, 0x02, false)
	cpu.add(vIndex, 0x03, false)

	expected := 0x05

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x0 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
	}
}

func TestADDWithCarry(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1
	cpu.v[vIndex] = 0xEE

	cpu.add(vIndex, 0xEE, true)

	expected := 0xEE + 0xEE

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x1 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x1)
	}
}

func TestSUBWithoutCarry(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1

	cpu.sub(vIndex, 0x02, 0x03)

	expected := 0x02 - 0x03

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x0 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
	}
}

func TestSUBWithCarry(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1

	cpu.sub(vIndex, 0x03, 0x02)

	expected := 0x03 - 0x02

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x1 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x1)
	}
}

func TestADIWithoutOverflow(t *testing.T) {
	cpu := NewCpu()

	cpu.adi(0x80)
	cpu.adi(0x50)

	var expected uint16 = 0x80 + 0x50

	if cpu.i != expected {
		t.Errorf("cpu.i = 0x%X; expected 0x%X", cpu.i, expected)
	}

	if cpu.v[0xF] != 0x0 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
	}
}

func TestADIWithOverflow(t *testing.T) {
	cpu := NewCpu()

	cpu.adi(0x0FFF)
	cpu.adi(0x01)

	var expected uint16 = 0x0FFF + 0x01

	if cpu.i != expected {
		t.Errorf("cpu.i = 0x%X; expected 0x%X", cpu.i, expected)
	}

	if cpu.v[0xF] != 0x1 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
	}
}

func TestSHRWithoutFlag(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1
	cpu.v[vIndex] = 0b11111110
	expected := cpu.v[vIndex] >> 1

	cpu.shr(vIndex)

	if cpu.v[vIndex] != expected {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x0 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
	}
}

func TestSHRWithFlag(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1
	cpu.v[vIndex] = 0b00000001
	expected := cpu.v[vIndex] >> 1

	cpu.shr(vIndex)

	if cpu.v[vIndex] != expected {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x1 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x1)
	}
}

func TestSHLWithoutFlag(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1
	cpu.v[vIndex] = 0b01111110
	expected := cpu.v[vIndex] << 1

	cpu.shl(vIndex)

	if cpu.v[vIndex] != expected {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x0 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
	}
}

func TestSHLWithFlag(t *testing.T) {
	cpu := NewCpu()

	var vIndex uint8 = 0x1
	cpu.v[vIndex] = 0b11111110
	expected := cpu.v[vIndex] << 1

	cpu.shl(vIndex)

	if cpu.v[vIndex] != expected {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x1 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x0)
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

func TestDRWNoWrapAndNoCollision(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0xD3)
	cpu.mmu.Write(0x201, 0xD2)
	cpu.i = 0x300
	cpu.v[0x3] = 0
	cpu.v[0xD] = 0
	cpu.mmu.Write(0x300, 0x11)
	cpu.mmu.Write(0x301, 0x88)

	cpu.cls()

	cpu.Clock()

	if cpu.v[0xF] != 0x00 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x00", cpu.v[0xF])
	}

	if cpu.display[0][3] != 0x01 {
		t.Errorf("cpu.display[0][3] = 0x%X; expected 0x01", cpu.display[0][3])
	}

	if cpu.display[0][7] != 0x01 {
		t.Errorf("cpu.display[0][7] = 0x%X; expected 0x01", cpu.display[0][7])
	}

	if cpu.display[1][0] != 0x01 {
		t.Errorf("cpu.display[1][0] = 0x%X; expected 0x01", cpu.display[1][0])
	}

	if cpu.display[1][4] != 0x01 {
		t.Errorf("cpu.display[1][4] = 0x%X; expected 0x01", cpu.display[1][4])
	}
}

func TestSKP(t *testing.T) {
	cpu := NewCpu()

	cpu.skp(false)

	if cpu.pc != 0x200 {
		t.Errorf("cpu.pc = 0x%X; expected 0x200", cpu.pc)
	}

	cpu.skp(true)

	if cpu.pc != 0x202 {
		t.Errorf("cpu.pc = 0x%X; expected 0x202", cpu.pc)
	}
}
