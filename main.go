package main

func main() {
	cpu := new(CPU)

	cpu.pc = 0x200 // start address

	cpu.LoadFont()

	cpu.Clock()
}
