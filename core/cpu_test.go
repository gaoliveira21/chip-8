package core

import (
	"os"
	"slices"
	"testing"
)

func TestNewCpu(t *testing.T) {
	cpu := NewCpu()

	inMemoryFonts := [len(FontData)]byte{}

	for i := 0x050; i <= 0x09F; i++ {
		font := cpu.mmu.Fetch(uint16(i))

		inMemoryFonts[i-0x050] = byte(font >> 8)
	}

	if cpu.pc != 0x200 {
		t.Errorf("cpu.pc = %d; expected 0x200", cpu.pc)
	}

	if inMemoryFonts != FontData {
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
	romData, err := os.ReadFile("../cli/roms/IBM.ch8")

	if err != nil {
		t.Fatal(err)
	}

	cpu := NewCpu()

	cpu.LoadROM(romData)

	inMemoryROM := []byte{}

	for i := 0; i < len(romData); i++ {
		romByte := cpu.mmu.Fetch(uint16(i + 0x200))

		inMemoryROM = append(inMemoryROM, byte(romByte>>8))
	}

	if !slices.Equal[[]byte](inMemoryROM, romData) {
		t.Error("Error loading ROM")
	}
}

func TestCLS(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x00)
	cpu.mmu.Write(0x201, 0xE0)

	cpu.Display[0][0] = 0xFF
	cpu.Display[0][1] = 0xEF

	cpu.Clock()

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			if cpu.Display[i][j] != 0x00 {
				t.Errorf("cpu.Display[%d][%d] = 0x%X; expected 0x00", i, j, cpu.Display[i][j])
			}
		}
	}
}

func TestRET(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x00)
	cpu.mmu.Write(0x201, 0xEE)

	cpu.mmu.Stack.Push(0xDDEE)
	cpu.Clock()

	if cpu.pc != 0xDDEE {
		t.Errorf("cpu.pc = 0x%X; expected 0xDDEE", cpu.pc)
	}
}

func TestJP0x0000(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x01)
	cpu.mmu.Write(0x201, 0x11)

	cpu.Clock()

	if cpu.pc != 0x0111 {
		t.Errorf("cpu.pc = 0x%X; expected 0x0111", cpu.pc)
	}
}

func TestJP0x1000(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x12)
	cpu.mmu.Write(0x201, 0x34)

	cpu.Clock()

	if cpu.pc != 0x0234 {
		t.Errorf("cpu.pc = 0x%X; expected 0x0111", cpu.pc)
	}
}

func TestCALL(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x300, 0x24)
	cpu.mmu.Write(0x301, 0x00)

	cpu.pc = 0x300
	cpu.Clock()

	stackPC := cpu.mmu.Stack.Pop()
	currentPC := cpu.pc

	if stackPC != 0x302 {
		t.Errorf("Stack PC = 0x%X; expected 0x300", stackPC)
	}

	if currentPC != 0x400 {
		t.Errorf("Current PC = 0x%X; expected 0x400", currentPC)
	}
}

func TestJPWithOffset(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0xBF)
	cpu.mmu.Write(0x201, 0xF0)

	cpu.v[0x0] = 0x02

	expected := 0xFF0 + uint16(cpu.v[0x0])

	cpu.Clock()

	if cpu.pc != uint16(expected) {
		t.Errorf("cpu.pc = 0x%X; expected 0x%X", cpu.pc, expected)
	}
}

func TestSKPVxEqualsToNN(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x35)
	cpu.mmu.Write(0x201, 0x68)

	cpu.v[0x5] = 0x68

	cpu.Clock()

	if cpu.pc != 0x204 {
		t.Errorf("cpu.pc = 0x%X; expected 0x204", cpu.pc)
	}

	cpu.mmu.Write(0x204, 0x35)
	cpu.mmu.Write(0x205, 0x70)

	cpu.Clock()

	if cpu.pc != 0x206 {
		t.Errorf("cpu.pc = 0x%X; expected 0x206", cpu.pc)
	}
}

func TestSKPVxNotEqualToNN(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x45)
	cpu.mmu.Write(0x201, 0x70)

	cpu.v[0x5] = 0x68

	cpu.Clock()

	if cpu.pc != 0x204 {
		t.Errorf("cpu.pc = 0x%X; expected 0x204", cpu.pc)
	}

	cpu.mmu.Write(0x204, 0x45)
	cpu.mmu.Write(0x205, 0x68)

	cpu.Clock()

	if cpu.pc != 0x206 {
		t.Errorf("cpu.pc = 0x%X; expected 0x206", cpu.pc)
	}
}

func TestSKPVxEqualsToVy(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x55)
	cpu.mmu.Write(0x201, 0x60)

	cpu.v[0x5] = 0x68
	cpu.v[0x6] = 0x68

	cpu.Clock()

	if cpu.pc != 0x204 {
		t.Errorf("cpu.pc = 0x%X; expected 0x204", cpu.pc)
	}

	cpu.mmu.Write(0x204, 0x55)
	cpu.mmu.Write(0x205, 0x60)

	cpu.v[0x5] = 0x70
	cpu.v[0x6] = 0x71

	cpu.Clock()

	if cpu.pc != 0x206 {
		t.Errorf("cpu.pc = 0x%X; expected 0x206", cpu.pc)
	}
}

func TestLDNNToVx(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x61)
	cpu.mmu.Write(0x201, 0xFF)

	cpu.Clock()

	if cpu.v[0x1] != 0xFF {
		t.Errorf("cpu.v[0x1] = 0x%X; expected 0x%X", cpu.v[0x1], 0xFF)
	}
}

func TestLDVyToVx(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x81)
	cpu.mmu.Write(0x201, 0x20)

	cpu.v[0x2] = 0x60

	cpu.Clock()

	if cpu.v[0x1] != 0x60 {
		t.Errorf("cpu.v[0x1] = 0x%X; expected 0x%X", cpu.v[0x1], 0x60)
	}
}

func TestADD(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x71)
	cpu.mmu.Write(0x201, 0x03)

	var vIndex uint8 = 0x1

	cpu.v[vIndex] = 0x02

	cpu.Clock()

	expected := 0x05

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}
}

func TestADDWitoutCarry(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x81)
	cpu.mmu.Write(0x201, 0x24)

	var vIndex uint8 = 0x1
	var yIndex uint8 = 0x2

	cpu.v[vIndex] = 0x02
	cpu.v[yIndex] = 0x03

	cpu.Clock()

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

	cpu.mmu.Write(0x200, 0x81)
	cpu.mmu.Write(0x201, 0x24)

	var vIndex uint8 = 0x1
	var yIndex uint8 = 0x2

	cpu.v[vIndex] = 0xEE
	cpu.v[yIndex] = 0xEE

	cpu.Clock()

	expected := 0xEE + 0xEE

	if cpu.v[vIndex] != byte(expected) {
		t.Errorf("cpu.v[%d] = 0x%X; expected 0x%X", vIndex, cpu.v[vIndex], expected)
	}

	if cpu.v[0xF] != 0x1 {
		t.Errorf("cpu.v[0xF] = 0x%X; expected 0x%X", cpu.v[0xF], 0x1)
	}
}

func TestVxORVy(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x81)
	cpu.mmu.Write(0x201, 0x21)

	cpu.v[0x1] = 0x05
	cpu.v[0x2] = 0x10

	expected := byte(0x05 | 0x10)

	cpu.Clock()

	if cpu.v[0x1] != expected {
		t.Errorf("cpu.v[0x1] = 0x%X; expected 0x%X", cpu.v[0x1], expected)
	}
}

func TestVxANDVy(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x81)
	cpu.mmu.Write(0x201, 0x22)

	cpu.v[0x1] = 0x05
	cpu.v[0x2] = 0x10

	expected := byte(0x05 & 0x10)

	cpu.Clock()

	if cpu.v[0x1] != expected {
		t.Errorf("cpu.v[0x1] = 0x%X; expected 0x%X", cpu.v[0x1], expected)
	}
}

func TestVxXORVy(t *testing.T) {
	cpu := NewCpu()

	cpu.mmu.Write(0x200, 0x81)
	cpu.mmu.Write(0x201, 0x23)

	cpu.v[0x1] = 0x05
	cpu.v[0x2] = 0x10

	expected := byte(0x05 ^ 0x10)

	cpu.Clock()

	if cpu.v[0x1] != expected {
		t.Errorf("cpu.v[0x1] = 0x%X; expected 0x%X", cpu.v[0x1], expected)
	}
}

// TODO: Implement with real memory addresses
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

	if cpu.SoundTimer != 0x80 {
		t.Errorf("cpu.SoundTimer = 0x%X; expected 0x80", cpu.SoundTimer)
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

func TestBCD(t *testing.T) {
	cpu := NewCpu()
	cpu.i = 0x300

	cpu.bcd(128)

	v := uint8(cpu.mmu.Fetch(cpu.i) >> 8)
	if v != 1 {
		t.Errorf("value at memory addr = 0x%X; expected 0x%X", cpu.i, 1)
	}

	v = uint8(cpu.mmu.Fetch(cpu.i+1) >> 8)
	if v != 2 {
		t.Errorf("value at memory addr = 0x%X; expected 0x%X", cpu.i, 2)
	}

	v = uint8(cpu.mmu.Fetch(cpu.i+2) >> 8)
	if v != 8 {
		t.Errorf("value at memory addr = 0x%X; expected 0x%X", cpu.i, 8)
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

	if cpu.Display[0][3] != 0x01 {
		t.Errorf("cpu.Display[0][3] = 0x%X; expected 0x01", cpu.Display[0][3])
	}

	if cpu.Display[0][7] != 0x01 {
		t.Errorf("cpu.Display[0][7] = 0x%X; expected 0x01", cpu.Display[0][7])
	}

	if cpu.Display[1][0] != 0x01 {
		t.Errorf("cpu.Display[1][0] = 0x%X; expected 0x01", cpu.Display[1][0])
	}

	if cpu.Display[1][4] != 0x01 {
		t.Errorf("cpu.Display[1][4] = 0x%X; expected 0x01", cpu.Display[1][4])
	}
}
