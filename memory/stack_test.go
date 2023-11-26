package memory_test

import (
	"testing"

	"github.com/gaoliveira21/chip8/memory"
)

func TestStack(t *testing.T) {
	stack := new(memory.Stack)

	stack.Push(0xFF)

	if stack.SP != 0x01 {
		t.Errorf("Stack.SP = %d; expected 0x01", stack.SP)
	}

	stack.Push(0xAA)
	stack.Push(0xBB)

	if stack.SP != 0x03 {
		t.Errorf("Stack.SP = %d; expected 0x03", stack.SP)
	}

	stackData := stack.Pop()

	if stack.SP != 0x02 {
		t.Errorf("Stack.SP = %d; expected 0x02", stack.SP)
	}

	if stackData != 0xBB {
		t.Errorf("Stack.Pop() = %d; expected 0xBB", stackData)
	}
}
