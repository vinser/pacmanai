package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vinser/pacmanai/internal/entity"
	"github.com/vinser/pacmanai/internal/maze"
	"github.com/vinser/pacmanai/internal/render"
	"github.com/vinser/pacmanai/internal/state"
)

type GameState int

const (
	StatePlaying GameState = iota
	StateRespawning
	StateGameOver
)

const (
	frightenedPeriod = 10 * time.Second
	respawnPeriod    = 3 * time.Second
)

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
	state             GameState
	respawnUntil      time.Time
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
		pacman: entity.NewPacman(entity.Position{X: 1, Y: 1}),
		ghosts: []*entity.Ghost{
			entity.NewGhost(entity.Blinky, entity.Position{X: 9, Y: 3}),
			entity.NewGhost(entity.Inky, entity.Position{X: 10, Y: 3}),
			entity.NewGhost(entity.Pinky, entity.Position{X: 9, Y: 5}),
			entity.NewGhost(entity.Clyde, entity.Position{X: 10, Y: 5}),
		},
		ghostTickInterval: 500 * time.Millisecond,
		lastGhostMove:     time.Now(),
		state:             StatePlaying,
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
	if m.state == StateRespawning {
		if time.Now().After(m.respawnUntil) {
			m.pacman.SetPos(m.pacman.Home())
			for _, g := range m.ghosts {
				g.SetPos(g.Home())
				g.SetState(entity.Chase)
			}
			m.state = StatePlaying
		}
		return m, tickGhosts()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state != StatePlaying {
			// Ignore input when Respawning or Game Over
			return m, nil
		}
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
				m.score.AddGhostPoints()
				g.SetState(entity.Eaten)
				g.SetPos(g.Home())
				return false
			case entity.Chase, entity.Scatter:
				m.pacman.LoseLife()
				if m.pacman.IsDead() {
					m.state = StateGameOver
					return true
				}
				// Enter respawn mode
				m.state = StateRespawning
				m.respawnUntil = time.Now().Add(respawnPeriod)
				return false
			}
		}
	}
	return false
}

// View renders the current game state.
func (m Model) View() string {
	switch m.state {
	case StateGameOver:
		return render.RenderGameOver(m.score)
	case StateRespawning:
		return render.RenderRespawning(m.pacman.Lives())
	default:
		return render.RenderAll(m.maze, m.pacman, m.ghosts, m.score)
	}
}
