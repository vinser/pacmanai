package maze

import (
	"errors"
	"strconv"
	"strings"
)

// Tile represents a type of cell in the maze.
type Tile int

const (
	Wall Tile = iota
	Dot
	Empty
	PowerPellet
)

// Maze represents the layout of the game field.
type Maze struct {
	width  int
	height int
	grid   [][]Tile
}

// Width returns the width of the maze.
func (m *Maze) Width() int {
	return m.width
}

// Height returns the height of the maze.
func (m *Maze) Height() int {
	return m.height
}

// TileAt returns the tile at the specified coordinates.
func (m *Maze) TileAt(x, y int) (Tile, error) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return Empty, errors.New("out of bounds")
	}
	return m.grid[y][x], nil
}

// SetTile sets a tile at the given coordinates.
func (m *Maze) SetTile(x, y int, tile Tile) error {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return errors.New("out of bounds")
	}
	m.grid[y][x] = tile
	return nil
}

// EatItem replaces a dot or power pellet with empty space and returns the eaten tile type.
func (m *Maze) EatItem(x, y int) Tile {
	tile := m.grid[y][x]
	if tile == Dot || tile == PowerPellet {
		m.grid[y][x] = Empty
	}
	return tile
}

// LoadDefault returns a static hardcoded maze for testing/demo.
func LoadDefault() *Maze {
	layout := []string{
		"####################",
		"#........##........#",
		"#.####.#.##.####.#.#",
		" o#  #.#.##.#  #.#o ",
		"#.####.#.##.####.#.#",
		"#..................#",
		"####################",
	}

	// Sanity check: tunnel sides must be symmetric (either both open or both closed)
	for y, row := range layout {
		if len(row) != len(layout[0]) {
			panic("invalid maze: inconsistent row widths")
		}
		left := row[0]
		right := row[len(row)-1]
		if (left == ' ' && right == '#') || (left == '#' && right == ' ') {
			panic("invalid maze: asymmetric tunnel on row " + strconv.Itoa(y))
		}
	}

	height := len(layout)
	width := len(layout[0])
	grid := make([][]Tile, height)

	for y, line := range layout {
		grid[y] = make([]Tile, width)
		for x, char := range strings.Split(line, "") {
			switch char {
			case "#":
				grid[y][x] = Wall
			case ".":
				grid[y][x] = Dot
			case "o":
				grid[y][x] = PowerPellet
			default:
				grid[y][x] = Empty
			}
		}
	}

	return &Maze{
		width:  width,
		height: height,
		grid:   grid,
	}
}

// IsTunnelRow returns true if row y has open sides (tunnel).
func (m *Maze) IsTunnelRow(y int) bool {
	return y >= 0 && y < m.height && m.grid[y][0] != Wall && m.grid[y][m.width-1] != Wall
}
