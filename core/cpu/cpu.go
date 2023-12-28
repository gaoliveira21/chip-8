package cpu

import (
	"math/rand"
	"time"

	"github.com/gaoliveira21/chip8/core/font"
	"github.com/gaoliveira21/chip8/core/graphics"
	"github.com/gaoliveira21/chip8/core/memory"
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

type CPU struct {
	pc         uint16   // Program Counter
	i          uint16   // I Register
	v          [16]byte // Variable registers
	mmu        memory.MMU
	Graphics   *graphics.Graphics
	delayTimer uint8
	SoundTimer uint8
	Keys       [16]uint8
}

func NewCpu() CPU {
	cpu := CPU{
		pc:       0x200,
		Graphics: graphics.NewGraphics(),
	}

	cpu.loadFont()

	return cpu
}

func (cpu *CPU) LoadROM(rom []byte) {
	for index, b := range rom {
		cpu.mmu.Write(uint16(index)+0x200, b)
	}
}

func (cpu *CPU) Run() {
	cpu.clock()

	if cpu.delayTimer > 0 {
		cpu.delayTimer--
	}

	if cpu.SoundTimer > 0 {
		cpu.SoundTimer--
	}
}

func (cpu *CPU) loadFont() {
	for i := 0x050; i <= 0x09F; i++ {
		cpu.mmu.Write(uint16(i), font.FontData[i-0x050])
	}
}

func (cpu *CPU) decode(data uint16) (oc *opcode) {
	return NewOpcode(data)
}

func (cpu *CPU) clock() {
	data := cpu.mmu.Fetch(cpu.pc)
	cpu.pc += 2

	opcode := cpu.decode(data)

	switch opcode.instruction {
	case 0x0000:
		switch opcode.nnn {
		case 0x0E0:
			cpu.cls()

		case 0x0EE:
			cpu.ret()

		default:
			cpu.jp(opcode.nnn, 0)
		}
	case 0x1000:
		cpu.jp(opcode.nnn, 0)
	case 0x2000:
		cpu.call(opcode.nnn)
	case 0x3000:
		cpu.skp(cpu.v[opcode.registerX] == opcode.nn)
	case 0x4000:
		cpu.skp(cpu.v[opcode.registerX] != opcode.nn)
	case 0x5000:
		cpu.skp(cpu.v[opcode.registerX] == cpu.v[opcode.registerY])
	case 0x6000:
		cpu.ld(opcode.registerX, opcode.nn)
	case 0x7000:
		cpu.add(opcode.registerX, opcode.nn, false)
	case 0x8000:
		switch opcode.n {
		case 0x0:
			cpu.ld(opcode.registerX, cpu.v[opcode.registerY])
		case 0x1:
			cpu.or(opcode.registerX, cpu.v[opcode.registerY])
		case 0x2:
			cpu.and(opcode.registerX, cpu.v[opcode.registerY])
		case 0x3:
			cpu.xor(opcode.registerX, cpu.v[opcode.registerY])
		case 0x4:
			cpu.add(opcode.registerX, cpu.v[opcode.registerY], true)
		case 0x5:
			cpu.sub(opcode.registerX, cpu.v[opcode.registerX], cpu.v[opcode.registerY])
		case 0x6:
			cpu.shr(opcode.registerX)
		case 0x7:
			cpu.sub(opcode.registerX, cpu.v[opcode.registerY], cpu.v[opcode.registerX])
		case 0xE:
			cpu.shl(opcode.registerX)
		}
	case 0x9000:
		cpu.skp(cpu.v[opcode.registerX] != cpu.v[opcode.registerY])
	case 0xA000:
		cpu.ldi(opcode.nnn)
	case 0xB000:
		cpu.jp(opcode.nnn, cpu.v[0x0])
	case 0xC000:
		cpu.rnd(opcode.registerX, opcode.nn)
	case 0xD000:
		cpu.drw(opcode)
	case 0xE000:
		switch opcode.nn {
		case 0x9E:
			cpu.skp(cpu.Keys[cpu.v[opcode.registerX]] == 0x01)
		case 0xA1:
			cpu.skp(cpu.Keys[cpu.v[opcode.registerX]] == 0x00)
		}
	case 0xF000:
		switch opcode.nn {
		case 0x07:
			cpu.ld(opcode.registerX, cpu.delayTimer)
		case 0x15:
			cpu.ldt(cpu.v[opcode.registerX])
		case 0x18:
			cpu.lds(cpu.v[opcode.registerX])
		case 0x0A:
			cpu.ldk(opcode.registerX)
		case 0x1E:
			cpu.adi(uint16(cpu.v[opcode.registerX]))
		case 0x29:
			cpu.ldi(0x050 + 5*uint16(cpu.v[opcode.registerX]))
		case 0x33:
			cpu.bcd(cpu.v[opcode.registerX])
		case 0x55:
			cpu.stm(opcode.registerX)
		case 0x65:
			cpu.ldm(opcode.registerX)
		}
	}
}

func (cpu *CPU) cls() {
	cpu.Graphics.Clear()
}

func (cpu *CPU) ret() {
	cpu.pc = cpu.mmu.Stack.Pop()
}

func (cpu *CPU) jp(addr uint16, offset uint8) {
	cpu.pc = addr + uint16(offset)
}

func (cpu *CPU) call(addr uint16) {
	cpu.mmu.Stack.Push(cpu.pc)
	cpu.pc = addr
}

func (cpu *CPU) skp(condition bool) {
	if condition {
		cpu.pc += 2
	}
}

func (cpu *CPU) ld(vIndex uint8, b byte) {
	cpu.v[vIndex] = b
}

func (cpu *CPU) add(vIndex uint8, b byte, carry bool) {
	if carry {
		cpu.v[0xF] = 0x0
		result := uint16(cpu.v[vIndex]) + uint16(b)

		if result > 255 {
			cpu.v[0xF] = 0x1
		}

		cpu.v[vIndex] = byte(result)
	} else {
		cpu.v[vIndex] += b
	}
}

func (cpu *CPU) sub(vIndex uint8, minuend byte, subtrahend byte) {
	if minuend > subtrahend {
		cpu.v[0xF] = 0x1
	} else {
		cpu.v[0xF] = 0x0
	}

	cpu.v[vIndex] = minuend - subtrahend
}

func (cpu *CPU) ldi(value uint16) {
	cpu.i = value
}

func (cpu *CPU) ldk(vIndex uint8) {
	for i, v := range cpu.Keys {
		if v == 0x01 {
			cpu.v[vIndex] = uint8(i)
			return
		}
	}

	cpu.pc -= 2
}

func (cpu *CPU) ldt(value uint8) {
	cpu.delayTimer = value
}

func (cpu *CPU) lds(value uint8) {
	cpu.SoundTimer = value
}

func (cpu *CPU) bcd(value uint8) {
	cpu.mmu.Write(cpu.i, (value/100)%10)
	cpu.mmu.Write(cpu.i+1, (value/10)%10)
	cpu.mmu.Write(cpu.i+2, value%10)
}

func (cpu *CPU) or(vIndex uint8, b byte) {
	cpu.v[vIndex] |= b
}

func (cpu *CPU) and(vIndex uint8, b byte) {
	cpu.v[vIndex] &= b
}

func (cpu *CPU) xor(vIndex uint8, b byte) {
	cpu.v[vIndex] ^= b
}

func (cpu *CPU) shr(vIndex uint8) {
	cpu.v[0xF] = cpu.v[vIndex] & 0x01

	cpu.v[vIndex] >>= 1
}

func (cpu *CPU) shl(vIndex uint8) {
	cpu.v[0xF] = (cpu.v[vIndex] & 0x80) >> 7

	cpu.v[vIndex] <<= 1
}

func (cpu *CPU) rnd(vIndex uint8, b byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomByte := byte(r.Intn(256))

	cpu.v[vIndex] = randomByte & b
}

func (cpu *CPU) adi(value uint16) {
	cpu.i += value

	if cpu.i > 0x0FFF {
		cpu.v[0xF] = 0x1
	}
}

func (cpu *CPU) stm(vIndex uint8) {
	for i := 0; uint8(i) <= vIndex; i++ {
		cpu.mmu.Write(cpu.i+uint16(i), cpu.v[i])
	}
}

func (cpu *CPU) ldm(vIndex uint8) {
	for i := 0; uint8(i) <= vIndex; i++ {
		cpu.v[i] = byte(cpu.mmu.Fetch(cpu.i+uint16(i)) >> 8)
	}
}

func (cpu *CPU) drw(oc *opcode) {
	x := cpu.v[oc.registerX] & 0x3F
	y := cpu.v[oc.registerY] & 0x1F
	cpu.v[0xF] = 0x00

	for i := 0; uint8(i) < oc.n; i++ {
		addr := uint16(i) + cpu.i
		pixels := byte(cpu.mmu.Fetch(addr) >> 8)
		xIndex := x

		for j := 0; j < 8; j++ {
			bit := (pixels >> (7 - j)) & 0x01
			pixelOnDisplay := cpu.Graphics.GetPixel(int(y), int(xIndex))

			if bit == 0x01 && pixelOnDisplay == 0x01 {
				cpu.v[0xF] = 0x01
			}

			cpu.Graphics.SetPixel(int(y), int(xIndex), pixelOnDisplay^bit)

			xIndex++

			if int(xIndex) >= cpu.Graphics.Width {
				break
			}
		}

		y++

		if int(y) >= cpu.Graphics.Height {
			break
		}
	}
}
