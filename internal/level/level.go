package level

import (
	"time"

	"github.com/vinser/pacmanai/internal/maze"
)

type Config struct {
	Index             int
	Maze              *maze.Maze
	RemainingDots     int
	GhostTickInterval time.Duration
}

// Create initializes a new level with its configuration and dot count.
func Create(index int) *Config {
	loader := maze.LoadDefault // TODO: Replace with level-specific loader if available
	m := loader()
	dotCount := countDots(m)

	return &Config{
		Index:             index,
		Maze:              m,
		RemainingDots:     dotCount,
		GhostTickInterval: ghostInterval(index),
	}
}

// ghostInterval returns the ghost movement interval based on level.
func ghostInterval(level int) time.Duration {
	base := 500 * time.Millisecond
	step := 25 * time.Millisecond
	calculated := base - time.Duration(level-1)*step

	if calculated < 100*time.Millisecond {
		return 100 * time.Millisecond
	}
	return calculated
}

// countDots scans the maze and returns the number of dot/power-pellet tiles.
func countDots(m *maze.Maze) int {
	count := 0
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			tile, _ := m.TileAt(x, y)
			if tile == maze.Dot || tile == maze.PowerPellet {
				count++
			}
		}
	}
	return count
}
