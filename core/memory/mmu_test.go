package memory_test

import (
	"testing"

	"github.com/gaoliveira21/chip8/core/memory"
)

func TestMMU(t *testing.T) {
	mmu := new(memory.MMU)

	mmu.Write(0x00, 0x00)
	mmu.Write(0x01, 0xEE)

	mmu.Write(0xFFE, 0xAA)
	mmu.Write(0xFFF, 0xBB)

	word1 := mmu.Fetch(0x00)
	word2 := mmu.Fetch(0xFFE)

	if word1 != 0x00EE {
		t.Errorf("Fetch(0x00) = %d; expected 0x00EE", word1)
	}

	if word2 != 0xAABB {
		t.Errorf("Fetch(0x00) = %d; expected 0xAABB", word2)
	}
}
