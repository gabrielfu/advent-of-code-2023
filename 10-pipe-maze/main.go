package main

import (
	"fmt"
	"math"
	"os"
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
	// fmt.Println("=====================================")
	fmt.Printf("pos=%c (%d, %d) ", tile, g.Pos.r, g.Pos.c)
	// fmt.Println()

	for _, mi := range AllMoveInstructions {
		// fmt.Printf("Considering moving %s\n", mi.Direction)

		if mi.Direction == g.Heading.Opposite() {
			// fmt.Printf("Going in reverse, skipping %s\n", mi.Direction)
			continue
		}

		if tile != 'S' && !IsConnected(mi.Direction.Opposite(), tile) {
			// fmt.Printf("Not connected, skipping %s\n", mi.Direction)
			continue
		}

		if g.Maze.MoveIsValid(mi, g.Pos) {
			fmt.Printf("Moving %s\n", mi.Direction)
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

func part2() {}

func main() {
	part1()
	part2()
}
