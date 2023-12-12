package main

import (
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
	Maze       []string
	Nr         int
	Nc         int
	Sr         int
	Sc         int
	LoopLength int
}

func NewMaze(lines []string) *Maze {
	nr := len(lines)
	nc := len(lines[0])

	// Find 'S' in lines
	sr := 0
	sc := 0
	for r := 0; r < nr; r++ {
		for c := 0; c < nc; c++ {
			if lines[r][c] == 'S' {
				sr = r
				sc = c
				break
			}
		}
	}

	return &Maze{
		Maze:       lines,
		Nr:         nr,
		Nc:         nc,
		Sr:         sr,
		Sc:         sc,
		LoopLength: 0,
	}
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

func Dfs(maze *Maze, r, c int, direction Direction, steps int) {
	// If we are at the end, return
	if steps > 0 && r == maze.Sr && c == maze.Sc {
		maze.LoopLength = steps
		return
	}

	if maze.LoopLength > 0 {
		return
	}

	if r < 0 || r >= maze.Nr || c < 0 || c >= maze.Nc {
		return
	}

	// If no pipe, return
	cur := maze.Maze[r][c]
	if cur == '.' {
		return
	}

	// validate direction is ok
	switch direction {
	case Up:
		if cur != '|' && cur != '7' && cur != 'F' {
			return
		}
	case Down:
		if cur != '|' && cur != 'J' && cur != 'L' {
			return
		}
	case Right:
		if cur != '-' && cur != 'J' && cur != '7' {
			return
		}
	case Left:
		if cur != '-' && cur != 'F' && cur != 'L' {
			return
		}
	default:
	}

	// traverse
	switch cur {
	case '|':
		if direction == Up {
			Dfs(maze, r-1, c, Up, steps+1)
		} else if direction == Down {
			Dfs(maze, r+1, c, Down, steps+1)
		}
	case '-':
		if direction == Left {
			Dfs(maze, r, c-1, Left, steps+1)
		} else if direction == Right {
			Dfs(maze, r, c+1, Right, steps+1)
		}
	case '7':
		if direction == Up {
			Dfs(maze, r, c-1, Left, steps+1)
		} else if direction == Right {
			Dfs(maze, r+1, c, Down, steps+1)
		}
	case 'F':
		if direction == Up {
			Dfs(maze, r, c+1, Right, steps+1)
		} else if direction == Left {
			Dfs(maze, r+1, c, Down, steps+1)
		}
	case 'J':
		if direction == Down {
			Dfs(maze, r, c-1, Left, steps+1)
		} else if direction == Right {
			Dfs(maze, r-1, c, Up, steps+1)
		}
	case 'L':
		if direction == Down {
			Dfs(maze, r, c+1, Right, steps+1)
		} else if direction == Left {
			Dfs(maze, r-1, c, Up, steps+1)
		}
	case 'S':
		Dfs(maze, r-1, c, Up, steps+1)
		Dfs(maze, r+1, c, Down, steps+1)
		Dfs(maze, r, c+1, Right, steps+1)
		Dfs(maze, r, c-1, Left, steps+1)
	default:
		return
	}
}

func part1() {
	lines := ReadLines("input.txt")

	maze := NewMaze(lines)
	Dfs(maze, maze.Sr, maze.Sc, None, 0)
	farthest := int(math.Ceil(float64(maze.LoopLength) / 2))
	println(farthest)
}

func part2() {}

func main() {
	part1()
	part2()
}
