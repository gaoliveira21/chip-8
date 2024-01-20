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
	Instruction uint16
	RegisterX   uint8
	RegisterY   uint8
	N           uint8
	NN          uint8
	NNN         uint16
}

func NewOpcode(data uint16) *opcode {
	return &opcode{
		Instruction: data & INSTRUCTION_BITMASK,
		RegisterX:   uint8((data & X_BITMASK) >> 8),
		RegisterY:   uint8((data & Y_BITMASK) >> 4),
		N:           uint8(data & N_BITMASK),
		NN:          uint8(data & NN_BITMASK),
		NNN:         (data & NNN_BITMASK),
	}
}
