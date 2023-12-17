package main

import (
	"flag"

	"github.com/gaoliveira21/chip8/core"
)

func main() {
	rom := flag.String("rom", "", "ROM path")
	flag.Parse()

	core.RunChip8(*rom)
}
