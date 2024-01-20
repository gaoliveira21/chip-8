package cpu

import (
	"testing"
)

func TestDecode(t *testing.T) {
	opcode := NewOpcode(0xABCD)

	var expected uint16 = 0xA000

	if opcode.Instruction != expected {
		t.Errorf("opcode.Instruction = 0x%X; expected 0x%X", opcode.Instruction, expected)
	}

	expected = 0xB
	if opcode.RegisterX != uint8(expected) {
		t.Errorf("opcode.RegisterX = 0x%X; expected 0x%X", opcode.RegisterX, expected)
	}

	expected = 0xC
	if opcode.RegisterY != uint8(expected) {
		t.Errorf("opcode.RegisterY = 0x%X; expected 0x%X", opcode.RegisterY, expected)
	}

	expected = 0xD
	if opcode.N != uint8(expected) {
		t.Errorf("opcode.N = 0x%X; expected 0x%X", opcode.N, expected)
	}

	expected = 0xCD
	if opcode.NN != uint8(expected) {
		t.Errorf("opcode.NN = 0x%X; expected 0x%X", opcode.NN, expected)
	}

	expected = 0xBCD
	if opcode.NNN != expected {
		t.Errorf("opcode.NNN = 0x%X; expected 0x%X", opcode.NNN, expected)
	}
}
