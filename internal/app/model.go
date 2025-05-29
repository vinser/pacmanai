package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vinser/pacmanai/internal/entity"
	"github.com/vinser/pacmanai/internal/level"
	"github.com/vinser/pacmanai/internal/maze"
	"github.com/vinser/pacmanai/internal/render"
	"github.com/vinser/pacmanai/internal/state"
)

type GameState int

const (
	StatePlaying GameState = iota
	StateRespawning
	StateGameOver
	StateLevelIntro
)

const (
	frightenedPeriod = 10 * time.Second
	respawnPeriod    = 3 * time.Second
	levelIntroPeriod = 3 * time.Second
)

// Model implements the bubbletea.Model interface.
type Model struct {
	level           *level.Config
	pacman          *entity.Pacman
	ghosts          []*entity.Ghost
	lastGhostMove   time.Time
	powerMode       bool
	powerModeUntil  time.Time
	score           *entity.Score
	state           GameState
	respawnUntil    time.Time
	gameOver        bool
	levelIntroUntil time.Time
}

// NewModel initializes the game model with maze, player, and ghosts.
func NewModel() Model {
	st := state.Load()
	s := entity.NewScore()
	s.SetHigh(st.HighScore)
	lvl := level.Create(1)
	ghosts := []*entity.Ghost{
		entity.NewGhost(entity.Blinky, entity.Position{X: 9, Y: 3}),
		entity.NewGhost(entity.Inky, entity.Position{X: 10, Y: 3}),
		entity.NewGhost(entity.Pinky, entity.Position{X: 9, Y: 5}),
		entity.NewGhost(entity.Clyde, entity.Position{X: 10, Y: 5}),
	}
	return Model{
		level:         lvl,
		pacman:        entity.NewPacman(entity.Position{X: 1, Y: 1}),
		ghosts:        ghosts,
		lastGhostMove: time.Now(),
		score:         s,
		state:         StatePlaying,
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
	if m.state == StateLevelIntro {
		if time.Now().After(m.levelIntroUntil) {
			m.state = StatePlaying
		}
		return m, tickGhosts()
	}
	if m.state == StateRespawning {
		if time.Now().After(m.respawnUntil) {
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
			m.pacman.Move(m.level.Maze)

			pos := m.pacman.Pos()
			tile := m.level.Maze.EatItem(pos.X, pos.Y)

			switch tile {
			case maze.Dot:
				m.score.Add(10)
				m.level.RemainingDots--
			case maze.PowerPellet:
				m.score.Add(50)
				m.level.RemainingDots--
				m.powerMode = true
				m.powerModeUntil = time.Now().Add(frightenedPeriod)
				for _, g := range m.ghosts {
					g.SetState(entity.Frightened)
				}
			}
			if m.level.RemainingDots < 1 {
				m.advanceLevel()
			}

		}
	}

	m.updatePowerMode()

	switch msg.(type) {
	case tickMsg:
		if time.Since(m.lastGhostMove) >= m.level.GhostTickInterval {
			entity.MoveGhosts(m.ghosts, m.level.Maze, m.powerMode)
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
				m.pacman.SetPos(m.pacman.Home())
				for _, g := range m.ghosts {
					g.SetPos(g.Home())
					g.SetState(entity.Chase)
				}
				m.state = StateRespawning
				m.respawnUntil = time.Now().Add(respawnPeriod)
				return false
			}
		}
	}
	return false
}

func (m *Model) advanceLevel() {
	m.level = level.Create(m.level.Index + 1)
	m.pacman.SetPos(m.pacman.Home())
	for _, g := range m.ghosts {
		g.SetPos(g.Home())
		g.SetState(entity.Chase)
	}
	m.state = StateLevelIntro
	m.levelIntroUntil = time.Now().Add(levelIntroPeriod)
}

// View renders the current game state.
func (m Model) View() string {
	switch m.state {
	case StateLevelIntro:
		return render.RenderLevelIntro(m.level.Index)
	case StateGameOver:
		return render.RenderGameOver(m.score)
	case StateRespawning:
		return render.RenderRespawning(m.pacman.Lives())
	default:
		return render.RenderAll(m.level.Maze, m.pacman, m.ghosts, m.score, m.level.Index)
	}
}
