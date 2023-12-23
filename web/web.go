package main

import (
	"log"
	"os"
	"path"
	"syscall/js"

	"github.com/gaoliveira21/chip8/core"
	"github.com/gaoliveira21/chip8/web/http"
)

func main() {
	os.Setenv("EXEC_MODE", "web")

	document := js.Global().Get("document")

	romName := document.Get("rom").String()

	rom, err := http.ReadFile(path.Join("roms", romName))

	if err != nil {
		log.Fatal(err)
	}

	core.RunChip8(rom, "[CHIP-8] - "+romName)
}
