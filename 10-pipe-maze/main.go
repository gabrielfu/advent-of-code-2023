package main

import (
	"math"
	"os"
	"slices"
	"strings"
)

func ReadLines(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

type Maze struct {
	Maze []string
	Nr   int
	Nc   int
}

func NewMaze(lines []string) *Maze {
	nr := len(lines)
	nc := len(lines[0])

	return &Maze{
		Maze: lines,
		Nr:   nr,
		Nc:   nc,
	}
}

func (m *Maze) FindStart() (int, int) {
	for r, row := range m.Maze {
		for c, col := range row {
			if col == 'S' {
				return r, c
			}
		}
	}
	return -1, -1
}

func (m *Maze) GetTile(r, c int) byte {
	return m.Maze[r][c]
}

func (m *Maze) SetTile(r, c int, tile byte) {
	m.Maze[r] = m.Maze[r][:c] + string(tile) + m.Maze[r][c+1:]
}

func IsConnected(direction Direction, tile byte) bool {
	switch direction {
	case Up:
		return tile == '|' || tile == '7' || tile == 'F'
	case Down:
		return tile == '|' || tile == 'J' || tile == 'L'
	case Right:
		return tile == '-' || tile == 'J' || tile == '7'
	case Left:
		return tile == '-' || tile == 'F' || tile == 'L'
	default:
		return false
	}
}

func (m *Maze) MoveIsValid(mi MoveInstruction, curPos Coord) bool {
	newPos := curPos.Move(mi)
	if newPos.r < 0 || newPos.r >= m.Nr || newPos.c < 0 || newPos.c >= m.Nc {
		return false
	}
	newTile := m.GetTile(newPos.r, newPos.c)
	return newTile == 'S' || IsConnected(mi.Direction, newTile)
}

type Direction int

const (
	None Direction = iota
	Up
	Down
	Right
	Left
)

func (d Direction) String() string {
	switch d {
	case None:
		return "None"
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Right:
		return "Right"
	case Left:
		return "Left"
	default:
		return "Unknown"
	}
}

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Right:
		return Left
	case Left:
		return Right
	default:
		return None
	}
}

type Coord struct {
	r int
	c int
}

type MoveInstruction struct {
	Delta     Coord
	Direction Direction
}

var AllMoveInstructions = []MoveInstruction{
	{Coord{-1, 0}, Up},
	{Coord{1, 0}, Down},
	{Coord{0, 1}, Right},
	{Coord{0, -1}, Left},
}

func (c Coord) Move(m MoveInstruction) Coord {
	return Coord{c.r + m.Delta.r, c.c + m.Delta.c}
}

// To returns the move instruction to go from c to other
func (c Coord) To(other Coord) MoveInstruction {
	if other.r == c.r {
		if other.c == c.c+1 {
			return MoveInstruction{Coord{0, 1}, Right}
		} else if other.c == c.c-1 {
			return MoveInstruction{Coord{0, -1}, Left}
		}
	} else if other.c == c.c {
		if other.r == c.r+1 {
			return MoveInstruction{Coord{1, 0}, Down}
		} else if other.r == c.r-1 {
			return MoveInstruction{Coord{-1, 0}, Up}
		}
	}
	panic("invalid move")
}

type Game struct {
	Maze    *Maze
	Pos     Coord
	Heading Direction
	Steps   int
}

func NewGame(maze *Maze) *Game {
	r, c := maze.FindStart()
	return &Game{
		Maze:    maze,
		Pos:     Coord{r, c},
		Heading: None,
		Steps:   0,
	}
}

func (g *Game) GetTile() byte {
	return g.Maze.GetTile(g.Pos.r, g.Pos.c)
}

func (g *Game) Move() {
	tile := g.GetTile()
	for _, mi := range AllMoveInstructions {
		if mi.Direction == g.Heading.Opposite() {
			continue
		}

		if tile != 'S' && !IsConnected(mi.Direction.Opposite(), tile) {
			continue
		}

		if g.Maze.MoveIsValid(mi, g.Pos) {
			g.Pos = g.Pos.Move(mi)
			g.Heading = mi.Direction
			g.Steps++
			return
		}
	}
	panic("No valid move")
}

func part1() {
	lines := ReadLines("input.txt")

	maze := NewMaze(lines)
	game := NewGame(maze)

	for {
		game.Move()
		if game.GetTile() == 'S' {
			break
		}
	}
	farthest := int(math.Ceil(float64(game.Steps) / 2))
	println(farthest)
}

func part2() {
	lines := ReadLines("input.txt")

	maze := NewMaze(lines)
	game := NewGame(maze)

	visited := make(map[Coord]struct{})
	var firstPos Coord
	var lastPos Coord
	for {
		game.Move()
		if firstPos == (Coord{}) {
			firstPos = game.Pos
		}
		visited[game.Pos] = struct{}{}
		if game.GetTile() == 'S' {
			break
		}
		lastPos = game.Pos
	}

	// substitute 'S' with the underlying pipe
	// current game.Pos is 'S'
	firstMoveDir := game.Pos.To(firstPos).Direction
	lastMoveDir := lastPos.To(game.Pos).Direction
	var tile byte
	if firstMoveDir == Up && lastMoveDir == Up {
		tile = '|'
	} else if firstMoveDir == Down && lastMoveDir == Down {
		tile = '|'
	} else if firstMoveDir == Right && lastMoveDir == Right {
		tile = '-'
	} else if firstMoveDir == Left && lastMoveDir == Left {
		tile = '-'
	} else if firstMoveDir == Up && lastMoveDir == Right {
		tile = 'J'
	} else if firstMoveDir == Up && lastMoveDir == Left {
		tile = 'L'
	} else if firstMoveDir == Down && lastMoveDir == Right {
		tile = '7'
	} else if firstMoveDir == Down && lastMoveDir == Left {
		tile = 'F'
	} else if firstMoveDir == Right && lastMoveDir == Up {
		tile = 'F'
	} else if firstMoveDir == Right && lastMoveDir == Down {
		tile = 'L'
	} else if firstMoveDir == Left && lastMoveDir == Up {
		tile = '7'
	} else if firstMoveDir == Left && lastMoveDir == Down {
		tile = 'J'
	} else {
		panic("invalid move")
	}
	maze.SetTile(game.Pos.r, game.Pos.c, tile)

	// point in polygon
	area := 0
	for r := 0; r < maze.Nr; r++ {
		for c := 0; c < maze.Nc; c++ {
			if _, ok := visited[Coord{r, c}]; ok {
				continue
			}

			crossings := 0
			for x := c + 1; x < maze.Nc; x++ {
				tile := maze.GetTile(r, x)
				_, ok := visited[Coord{r, x}]
				if ok && slices.Contains([]byte{'|', 'J', 'L'}, tile) {
					crossings++
				}
			}
			if crossings%2 == 1 {
				area++
			}
		}
	}
	println(area)

}

func main() {
	part1()
	part2()
}
