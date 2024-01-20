package core

import (
	"fmt"

	"github.com/gaoliveira21/chip8/core/cpu"
)

func Disassemble(bytes []byte) {
	fmt.Printf("|Bytes   |Op  |Nemonic\n")

	for i := 0; i < len(bytes); i += 2 {
		hb := uint16(bytes[i])
		lb := uint16(bytes[i+1])

		instruction := (hb << 8) | lb

		opcode := cpu.NewOpcode(instruction)

		switch opcode.Instruction {
		case 0x0000:
			if opcode.RegisterY == 0xC {
				fmt.Printf("0x%.4X - 00CN (SCROLL-DOWN N)\n", instruction)
				return
			}

			switch opcode.NNN {
			case 0x0E0:
				fmt.Printf("0x%.4X - 00E0 CLS\n", instruction)

			case 0x0EE:
				fmt.Printf("0x%.4X - 00EE RET\n", instruction)

			case 0x0FE:
				fmt.Printf("0x%.4X - 00FE (LORES)\n", instruction)

			case 0x0FF:
				fmt.Printf("0x%.4X - 00FF (HIRES)\n", instruction)

			case 0x0FB:
				fmt.Printf("0x%.4X - 00FB (SCROLL-RIGHT)\n", instruction)

			case 0x0FC:
				fmt.Printf("0x%.4X - 00FC (SCROLL-LEFT)\n", instruction)

			case 0x0FD:
				fmt.Printf("0x%.4X - 00FD (EXIT)\n", instruction)

			default:
				fmt.Printf("0x%.4X - 0NNN JP addr\n", instruction)
			}
		case 0x1000:
			fmt.Printf("0x%.4X - 1NNN JP addr\n", instruction)
		case 0x2000:
			fmt.Printf("0x%.4X - 2NNN CALL addr\n", instruction)
		case 0x3000:
			fmt.Printf("0x%.4X - 3XKK SE Vx, byte\n", instruction)
		case 0x4000:
			fmt.Printf("0x%.4X - 4XKK SNE Vx, byte\n", instruction)
		case 0x5000:
			fmt.Printf("0x%.4X - 5XY0 SE Vx, Vy\n", instruction)
		case 0x6000:
			fmt.Printf("0x%.4X - 6XKK LD Vx, byte\n", instruction)
		case 0x7000:
			fmt.Printf("0x%.4X - 7XKK ADD Vx, byte\n", instruction)
		case 0x8000:
			switch opcode.N {
			case 0x0:
				fmt.Printf("0x%.4X - 8XY0 LD Vx, Vy\n", instruction)
			case 0x1:
				fmt.Printf("0x%.4X - 8XY1 OR Vx, Vy\n", instruction)
			case 0x2:
				fmt.Printf("0x%.4X - 8XY2 AND Vx, Vy\n", instruction)
			case 0x3:
				fmt.Printf("0x%.4X - 8XY3 XOR Vx, Vy\n", instruction)
			case 0x4:
				fmt.Printf("0x%.4X - 8XY4 ADD Vx, Vy\n", instruction)
			case 0x5:
				fmt.Printf("0x%.4X - 8XY5 SUB Vx, Vy\n", instruction)
			case 0x6:
				fmt.Printf("0x%.4X - 8XY6 SHR Vx {, Vy}\n", instruction)
			case 0x7:
				fmt.Printf("0x%.4X - 8XY7 SUBN Vx, Vy\n", instruction)
			case 0xE:
				fmt.Printf("0x%.4X - 8XYE SHL Vx {, Vy}\n", instruction)
			}
		case 0x9000:
			fmt.Printf("0x%.4X - 9XY0 SNE Vx, Vy\n", instruction)
		case 0xA000:
			fmt.Printf("0x%.4X - ANNN LD I, addr\n", instruction)
		case 0xB000:
			fmt.Printf("0x%.4X - BNNN JP V0, addr\n", instruction)
		case 0xC000:
			fmt.Printf("0x%.4X - CXKK RND Vx, byte\n", instruction)
		case 0xD000:
			switch opcode.N {
			case 0x0:
				fmt.Printf("0x%.4X - DXY0 (SPRITE Vx Vy 0)\n", instruction)
			default:
				fmt.Printf("0x%.4X - DXYN DRW Vx, Vy, nibble\n", instruction)
			}
		case 0xE000:
			switch opcode.NN {
			case 0x9E:
				fmt.Printf("0x%.4X - EX9E SKP Vx\n", instruction)
			case 0xA1:
				fmt.Printf("0x%.4X - EXA1 SKNP Vx\n", instruction)
			}
		case 0xF000:
			switch opcode.NN {
			case 0x07:
				fmt.Printf("0x%.4X - FX07 LD Vx, DT\n", instruction)
			case 0x15:
				fmt.Printf("0x%.4X - FX15 LD DT, Vx\n", instruction)
			case 0x18:
				fmt.Printf("0x%.4X - FX18 LD ST, Vx\n", instruction)
			case 0x0A:
				fmt.Printf("0x%.4X - FX0A LD Vx, K\n", instruction)
			case 0x1E:
				fmt.Printf("0x%.4X - FX1E ADD I, Vx\n", instruction)
			case 0x29:
				fmt.Printf("0x%.4X - FX29 LD F, Vx\n", instruction)
			case 0x30:
				fmt.Printf("0x%.4X - FX30 (i := BIGHEX Vx)\n", instruction)
			case 0x33:
				fmt.Printf("0x%.4X - FX33 LD B, Vx\n", instruction)
			case 0x55:
				fmt.Printf("0x%.4X - FX55 LD [I], Vx\n", instruction)
			case 0x65:
				fmt.Printf("0x%.4X - FX65 LD Vx, [I]\n", instruction)
			case 0x75:
				fmt.Printf("0x%.4X - FX75 (SAVE FLAGS Vx)\n", instruction)
			case 0x85:
				fmt.Printf("0x%.4X - FX85 (LOAD FLAGS Vx)\n", instruction)
			}
		}
	}
}
