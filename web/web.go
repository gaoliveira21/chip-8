package main

import (
	"log"
	"os"

	"github.com/gaoliveira21/chip8/core"
	"github.com/gaoliveira21/chip8/web/http"
)

func main() {
	os.Setenv("EXEC_MODE", "web")

	rom, err := http.ReadFile("roms/PONG.ch8")

	if err != nil {
		log.Fatal(err)
	}

	core.RunChip8(rom, "[CHIP-8] - PONG")
}
