package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vinser/pacmanai/internal/entity"
	"github.com/vinser/pacmanai/internal/maze"
	"github.com/vinser/pacmanai/internal/render"
	"github.com/vinser/pacmanai/internal/state"
)

const frightenedPeriod = 10 * time.Second

// Model implements the bubbletea.Model interface.
type Model struct {
	score             *entity.Score
	maze              *maze.Maze
	pacman            *entity.Pacman
	ghosts            []*entity.Ghost
	ghostTickInterval time.Duration
	lastGhostMove     time.Time
	gameOver          bool
	powerMode         bool
	powerModeUntil    time.Time
}

// NewModel initializes the game model with maze, player, and ghosts.
func NewModel() Model {
	m := maze.LoadDefault()
	st := state.Load()
	s := entity.NewScore()
	s.SetHigh(st.HighScore)

	return Model{
		score:  s,
		maze:   m,
		pacman: entity.NewPacman(),
		ghosts: []*entity.Ghost{
			entity.NewGhost(entity.Blinky, entity.Position{X: 9, Y: 3}),
			entity.NewGhost(entity.Inky, entity.Position{X: 10, Y: 3}),
			entity.NewGhost(entity.Pinky, entity.Position{X: 9, Y: 5}),
			entity.NewGhost(entity.Clyde, entity.Position{X: 10, Y: 5}),
		},
		ghostTickInterval: 500 * time.Millisecond,
		lastGhostMove:     time.Now(),
	}
}

type tickMsg time.Time

func tickGhosts() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init is called once when the program starts.
func (m Model) Init() tea.Cmd {
	return tickGhosts()
}

// Update handles messages (e.g., key presses).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			m.pacman.HandleInput(msg)
			m.pacman.Move(m.maze)

			pos := m.pacman.Pos()
			tile := m.maze.EatItem(pos.X, pos.Y)

			switch tile {
			case maze.Dot:
				m.score.Add(10)
			case maze.PowerPellet:
				m.score.Add(50)
				m.powerMode = true
				m.powerModeUntil = time.Now().Add(frightenedPeriod)
				for _, g := range m.ghosts {
					g.SetState(entity.Frightened)
				}
			}
		}
	}

	m.updatePowerMode()

	switch msg.(type) {
	case tickMsg:
		if time.Since(m.lastGhostMove) >= m.ghostTickInterval {
			entity.MoveGhosts(m.ghosts, m.maze)
			m.lastGhostMove = time.Now()
		}
	}
	if m.checkCollisions() {
		m.gameOver = true
		currentScore := m.score.Get()
		st := state.Load()
		if currentScore > st.HighScore {
			st.HighScore = currentScore
			_ = state.Save(st)
		}
		return m, tea.Quit
	}

	return m, tickGhosts()
}

func (m *Model) updatePowerMode() {
	if m.powerMode && time.Now().After(m.powerModeUntil) {
		m.powerMode = false
		m.score.ResetGhostStreak()
		for _, g := range m.ghosts {
			if g.State() == entity.Frightened {
				g.SetState(entity.Chase)
			}
		}
	}
}

func (m *Model) checkCollisions() bool {
	pac := m.pacman.Pos()
	for _, g := range m.ghosts {
		if pac == g.Pos() {
			switch g.State() {
			case entity.Frightened:
				// Ghost is eaten - change state to Eaten, move to "home"
				m.score.AddGhostPoints()
				g.SetState(entity.Eaten)
				g.SetPos(g.Home())
			case entity.Chase, entity.Scatter:
				// Collision = Pac-Man death
				return true
			case entity.Eaten:
				// already eaten — ignore
			}
		}
	}
	return false
}

// View renders the current game state.
func (m Model) View() string {
	if m.gameOver {
		return render.RenderGameOver(m.score)
	}
	return render.RenderAll(m.maze, m.pacman, m.ghosts, m.score)
}
