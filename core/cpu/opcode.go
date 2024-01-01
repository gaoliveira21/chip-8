package cpu

const (
	INSTRUCTION_BITMASK = 0xF000
	X_BITMASK           = 0x0F00
	Y_BITMASK           = 0x00F0
	N_BITMASK           = 0x000F
	NN_BITMASK          = 0x00FF
	NNN_BITMASK         = 0x0FFF
)

type opcode struct {
	instruction uint16
	registerX   uint8
	registerY   uint8
	n           uint8
	nn          uint8
	nnn         uint16
}

func NewOpcode(data uint16) *opcode {
	return &opcode{
		instruction: data & INSTRUCTION_BITMASK,
		registerX:   uint8((data & X_BITMASK) >> 8),
		registerY:   uint8((data & Y_BITMASK) >> 4),
		n:           uint8(data & N_BITMASK),
		nn:          uint8(data & NN_BITMASK),
		nnn:         (data & NNN_BITMASK),
	}
}
