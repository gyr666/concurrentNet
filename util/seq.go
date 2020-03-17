package util

type Sequence struct {
	Max int
	now uint64
}

func (s *Sequence) Next() int {
	s.now++
	return int(s.now & uint64(s.Max))
}
