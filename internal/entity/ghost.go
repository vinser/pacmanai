package entity

import (
	"math/rand"

	"github.com/vinser/pacmanai/internal/maze"
)

// GhostState defines the current behavior mode of a ghost.
type GhostState int

const (
	Chase GhostState = iota
	Scatter
	Frightened
	Eaten
)

// GhostType defines the ghost's identity.
type GhostType int

const (
	Blinky GhostType = iota
	Inky
	Pinky
	Clyde
)

// Ghost represents a ghost entity.
type Ghost struct {
	position  Position
	direction Direction
	state     GhostState
	ghostType GhostType
}

// NewGhostWithType creates a ghost with specified type and position.
func NewGhostWithType(t GhostType, pos Position) *Ghost {
	return &Ghost{
		position:  pos,
		direction: Left,
		state:     Chase,
		ghostType: t,
	}
}

// Rune returns the character used to render the ghost.
func (g *Ghost) Rune() rune {
	switch g.ghostType {
	case Blinky:
		return 'B'
	case Inky:
		return 'I'
	case Pinky:
		return 'P'
	case Clyde:
		return 'Y'
	default:
		return 'G'
	}
}

// Pos returns the current position of the ghost.
func (g *Ghost) Pos() Position {
	return g.position
}

// SetPos sets the ghost's position directly.
func (g *Ghost) SetPos(pos Position) {
	g.position = pos
}

// State returns the current state of the ghost.
func (g *Ghost) State() GhostState {
	return g.state
}

// SetState updates the ghost's state.
func (g *Ghost) SetState(state GhostState) {
	g.state = state
}

// Move moves the ghost in its current direction.
// NextPos returns the position the ghost would move to.
func (g *Ghost) NextPos() Position {
	pos := g.position
	switch g.direction {
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

// Move tries to move the ghost forward if not hitting a wall.
func (g *Ghost) Move(m *maze.Maze) {
	next := g.NextPos()

	if m.IsTunnelRow(next.Y) {
		if next.X < 0 {
			next.X = m.Width() - 1
		} else if next.X >= m.Width() {
			next.X = 0
		}
	}

	tile, err := m.TileAt(next.X, next.Y)
	if err == nil && tile != maze.Wall {
		g.position = next
	}
}

// SetDirection sets the ghost's movement direction.
func (g *Ghost) SetDirection(dir Direction) {
	g.direction = dir
}

func (g *Ghost) MoveRandom(m *maze.Maze) {
	// Направления, кроме обратного
	possible := g.validDirectionsExcludingOpposite(m)

	// Если некуда идти кроме назад — разрешаем назад
	if len(possible) == 0 {
		possible = g.validAllDirections(m)
	}

	if len(possible) == 0 {
		return // полностью заблокирован
	}

	g.direction = possible[rand.Intn(len(possible))]
	g.Move(m)
}

func (g *Ghost) validDirectionsExcludingOpposite(m *maze.Maze) []Direction {
	var dirs []Direction
	opp := oppositeDirection(g.direction)

	for _, d := range []Direction{Up, Down, Left, Right} {
		if d == opp {
			continue
		}
		if canMoveTo(g.position, d, m) {
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func (g *Ghost) validAllDirections(m *maze.Maze) []Direction {
	var dirs []Direction
	for _, d := range []Direction{Up, Down, Left, Right} {
		if canMoveTo(g.position, d, m) {
			dirs = append(dirs, d)
		}
	}
	return dirs
}

// validDirections возвращает список допустимых направлений, кроме "назад".
func (g *Ghost) validDirections(m *maze.Maze) []Direction {
	var dirs []Direction
	opp := oppositeDirection(g.direction)

	for _, d := range []Direction{Up, Down, Left, Right} {
		if d == opp {
			continue // не разворачиваемся
		}
		pos := g.position.moveIn(d)

		if m.IsTunnelRow(pos.Y) {
			if pos.X < 0 {
				pos.X = m.Width() - 1
			} else if pos.X >= m.Width() {
				pos.X = 0
			}
		}

		tile, err := m.TileAt(pos.X, pos.Y)
		if err == nil && tile != maze.Wall {
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func canMoveTo(pos Position, d Direction, m *maze.Maze) bool {
	p := pos.moveIn(d)

	if m.IsTunnelRow(p.Y) {
		if p.X < 0 {
			p.X = m.Width() - 1
		} else if p.X >= m.Width() {
			p.X = 0
		}
	}

	tile, err := m.TileAt(p.X, p.Y)
	return err == nil && tile != maze.Wall
}

func (p Position) moveIn(d Direction) Position {
	switch d {
	case Up:
		p.Y--
	case Down:
		p.Y++
	case Left:
		p.X--
	case Right:
		p.X++
	}
	return p
}

func oppositeDirection(d Direction) Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	default:
		return d
	}
}
