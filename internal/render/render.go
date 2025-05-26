package render

import (
	"strings"

	"github.com/vinser/pacmanai/internal/entity"
	"github.com/vinser/pacmanai/internal/maze"
)

// isGhostAt returns rune of ghost if present at (x, y), else 0.
func ghostAt(x, y int, ghosts []*entity.Ghost) rune {
	for _, g := range ghosts {
		pos := g.Pos()
		if pos.X == x && pos.Y == y {
			return g.Rune()
		}
	}
	return 0
}

// RenderAll builds a full frame from the maze and entities.
func RenderAll(m *maze.Maze, pac *entity.Pacman, ghosts []*entity.Ghost) string {
	var sb strings.Builder

	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			// Ghosts first
			if ghost := ghostAt(x, y, ghosts); ghost != 0 {
				sb.WriteRune(ghost)
				continue
			}
			// Pacman
			if pac.Pos().X == x && pac.Pos().Y == y {
				sb.WriteRune('C')
				continue
			}
			// Maze tile
			tile, _ := m.TileAt(x, y)
			switch tile {
			case maze.Wall:
				sb.WriteRune('#')
			case maze.Dot:
				sb.WriteRune('.')
			case maze.PowerPellet:
				sb.WriteRune('o')
			default:
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
