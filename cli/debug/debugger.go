package debug

import (
	"fmt"
)

type DisassemblerResponse struct {
	Instructions []string `json:"instructions"`
	RomName      string   `json:"romName"`
}

func NewDebugger(rom []byte, romName string) {
	instructions := Disassemble(rom)

	for _, v := range instructions {
		fmt.Print(v)
	}
}
