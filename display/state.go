package display

type State struct {
	Cursor   int
	Max, Min int
}

func (s *State) Up() {
	if s.Cursor > 0 {
		s.Cursor--
	}
	if s.Cursor < s.Min {
		s.Min--
		s.Max--
	}
}

func (s *State) Down(listLen int) {
	if s.Cursor < listLen-1 {
		s.Cursor++
	}
	if s.Cursor > s.Max {
		s.Min++
		s.Max++
	}
}

func (s *State) NextPg(pgHeight int, listLen int) {
	s.Cursor += pgHeight
	if s.Cursor > listLen-1 {
		s.Cursor = listLen - 1
	}
	s.Min += pgHeight
	s.Max += pgHeight
	if s.Max >= listLen {
		s.Max = listLen - 1
		s.Min = s.Max - (pgHeight - 1)
	}
}

func (s *State) PrevPg(pgHeight int) {
	s.Cursor -= pgHeight
	if s.Cursor < 0 {
		s.Cursor = 0
	}
	s.Min -= pgHeight
	s.Max -= pgHeight
	if s.Min < 0 {
		s.Min = 0
		s.Max = pgHeight - 1
	}
}

func (s *State) Focus(i int, pgHeight int, listLen int) {
	s.Cursor = i
	s.Min = max(0, i-pgHeight/2)
	s.Max = min(listLen-1, s.Min+pgHeight-1)
}

func (s *State) SetMax(pgHeight int) {
	s.Max = max(s.Max, pgHeight-1)
}

func (s *State) Displayed(listLen int) int {
	return min(listLen, s.Max-s.Min+1)
}

func (s *State) IsSelected(n int) bool {
	return s.Cursor == n
}

func (s *State) Home(pgHeight int) {
	s.Cursor = 0
	s.Min = 0
	s.Max = pgHeight - 1
}

func (s *State) End(pgHeight int, listLen int) {
	s.Cursor = listLen - 1
	s.Max = listLen - 1
	s.Min = max(0, listLen-pgHeight)
}
