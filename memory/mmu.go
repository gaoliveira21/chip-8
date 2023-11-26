package memory

const RAM_SIZE = 4096 // 4 KB

type MMU struct {
	memory [RAM_SIZE]uint8
	Stack  Stack
}

func (m *MMU) Fetch(addr uint16) uint16 {
	hb := uint16(m.memory[addr])

	lb := uint16(m.memory[addr+1])

	return (hb << 8) | lb
}

func (m *MMU) Write(addr uint16, data byte) {
	m.memory[addr] = data
}
