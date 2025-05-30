package entity

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinser/pacmanai/internal/maze"
)

// Default lives count for Pacman.
const defaultLives = 3

// Direction represents movement direction.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Position represents coordinates on the map.
type Position struct {
	X, Y int
}

// Pacman represents the player character.
type Pacman struct {
	home      Position
	position  Position
	direction Direction
	lives     int
	// You can add more fields here, e.g., animation frame, lives, etc.
}

// NewPacman returns a new Pacman instance with default values.
func NewPacman(home Position) *Pacman {
	return &Pacman{
		home:      home,
		position:  home, // starting position
		direction: Right,
		lives:     defaultLives,
	}
}

// Home returns Pacman's home position.
func (p *Pacman) Home() Position {
	return p.home
}

// Pos returns Pacman's current position.
func (p *Pacman) Pos() Position {
	return p.position
}

// Dir returns Pacman's current direction.
func (p *Pacman) Dir() Direction {
	return p.direction
}

// SetPos sets Pacman's position explicitly.
func (p *Pacman) SetPos(pos Position) {
	p.position = pos
}

// NextPos returns the position Pacman would move to based on direction.
func (p *Pacman) NextPos() Position {
	pos := p.position
	switch p.direction {
	case Up:
		pos.Y--
	case Down:
		pos.Y++
	case Left:
		pos.X--
	case Right:
		pos.X++
	}
	return pos
}

// Move attempts to move Pacman if the next tile is not a wall.
func (p *Pacman) Move(m *maze.Maze) {
	next := p.NextPos()

	// Handle tunnel wrapping
	if m.IsTunnelRow(next.Y) {
		if next.X < 0 {
			next.X = m.Width() - 1
		} else if next.X >= m.Width() {
			next.X = 0
		}
	}

	tile, err := m.TileAt(next.X, next.Y)
	if err == nil && tile != maze.Wall {
		p.position = next
	}
}

// HandleInput updates Pacman's direction based on user input.
func (p *Pacman) HandleInput(msg tea.KeyMsg) {
	switch msg.String() {
	case "up", "w":
		p.direction = Up
	case "down", "s":
		p.direction = Down
	case "left", "a":
		p.direction = Left
	case "right", "d":
		p.direction = Right
	}
}

// Lives returns Pacman's remaining lives.
func (p *Pacman) Lives() int {
	return p.lives
}

// LoseLife reduces Pacman's lives by 1.
func (p *Pacman) LoseLife() {
	if p.lives > 0 {
		p.lives--
	}
}

// IsDead returns true if Pacman has no lives left.
func (p *Pacman) IsDead() bool {
	return p.lives <= 0
}

// ReseLives resets Pacman's lives to the default value.
func (p *Pacman) ResetLives() {
	p.lives = defaultLives
}

// AddLife adds a life to Pacman.
func (p *Pacman) AddLife() {
	p.lives++
}
