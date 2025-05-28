package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/vinser/pacmanai/internal/entity"
	"github.com/vinser/pacmanai/internal/maze"
)

var (
	styleFrightened = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	styleEaten      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	headerStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
)

func ghostAt(x, y int, ghosts []*entity.Ghost) *entity.Ghost {
	for _, g := range ghosts {
		if g.Pos().X == x && g.Pos().Y == y {
			return g
		}
	}
	return nil
}

// RenderAll returns the complete screen output with game entities and stats.
func RenderAll(m *maze.Maze, pac *entity.Pacman, ghosts []*entity.Ghost, score *entity.Score) string {
	var sb strings.Builder

	// Draw game header
	header := fmt.Sprintf("Score: %d   High Score: %d   Lives: %d\n", score.Get(), score.GetHigh(), pac.Lives())
	sb.WriteString(headerStyle.Render(header))
	sb.WriteRune('\n')

	// Draw maze with entities
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			if ghost := ghostAt(x, y, ghosts); ghost != nil {
				sb.WriteString(RenderGhost(ghost))
				continue
			}
			if pac.Pos().X == x && pac.Pos().Y == y {
				sb.WriteRune('C')
				continue
			}
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

	// Controls footer
	sb.WriteString("\nControls: ← ↑ ↓ → — move, q — quit\n")
	return sb.String()
}

func RenderGhost(g *entity.Ghost) string {
	switch g.State() {
	case entity.Frightened:
		return styleFrightened.Render("v")
	case entity.Eaten:
		return styleEaten.Render("x")
	default:
		return string(g.Rune())
	}
}

func RenderGameOver(score *entity.Score) string {
	var msg strings.Builder
	msg.WriteString("\nGame Over!\n")
	if score.Get() > score.GetHigh() {
		msg.WriteString("!!! New High Score: ")
	} else {
		msg.WriteString("Your Score: ")
	}
	msg.WriteString(fmt.Sprintf("%d\n", score.Get()))
	return msg.String()
}

func RenderRespawning(lives int) string {
	var flash string
	if (time.Now().UnixNano()/int64(time.Millisecond)/500)%2 == 0 {
		flash = "Respawning..."
	} else {
		flash = ""
	}
	return fmt.Sprintf(
		"\n%s\nLives: %d\nGet ready to continue\n",
		flash,
		lives,
	)
}
