package entity

type Score struct {
	value             int
	eatenGhostsStreak int
}

func NewScore() *Score {
	return &Score{}
}

func (s *Score) Add(points int) {
	s.value += points
}

func (s *Score) Get() int {
	return s.value
}

func (s *Score) Reset() {
	s.value = 0
	s.eatenGhostsStreak = 0
}

// Call when Pac-Man eats a frightened ghost
func (s *Score) AddGhostPoints() {
	points := 200 << s.eatenGhostsStreak // 200, 400, 800, 1600
	s.value += points
	if s.eatenGhostsStreak < 3 {
		s.eatenGhostsStreak++
	}
}

// Reset when frightened mode ends
func (s *Score) ResetGhostStreak() {
	s.eatenGhostsStreak = 0
}
