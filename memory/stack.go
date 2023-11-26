package memory

type Stack struct {
	data [16]uint16
	SP   uint16 // Stack Pointer
}

func (s *Stack) Push(addr uint16) {
	s.data[s.SP] = addr
	s.SP++
}

func (s *Stack) Pop() uint16 {
	s.SP--
	addr := s.data[s.SP]

	s.data[s.SP] = 0x00

	return addr
}
