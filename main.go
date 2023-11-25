package main

import "fmt"

func main() {
	cpu := new(CPU)

	cpu.pc = 0x200 // start address

	cpu.LoadFont()

	fmt.Printf("%d", cpu.memory)
}
