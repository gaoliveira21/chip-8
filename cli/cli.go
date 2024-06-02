package main

import (
	"flag"
	"log"
	"os"

	"github.com/gaoliveira21/chip8/cli/debug"
	"github.com/gaoliveira21/chip8/core"
)

func main() {
	rom := flag.String("rom", "", "ROM path")
	flag.Parse()

	romData, err := os.ReadFile(*rom)

	go debug.NewDebugger(romData, *rom)

	if err != nil {
		log.Fatal(err)
	}

	core.RunChip8(romData, *rom)
}
