package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadLines(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

type Direction string

const (
	Up    Direction = "U"
	Down  Direction = "D"
	Left  Direction = "L"
	Right Direction = "R"
)

type Instruction struct {
	Direction Direction
	Length    int
	Color     string
}

type Plan []Instruction

func NewPlan(lines []string) Plan {
	var plan Plan
	for _, line := range lines {
		split := strings.Split(line, " ")
		l, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}
		plan = append(plan, Instruction{
			Direction: Direction(split[0]),
			Length:    l,
			Color:     split[2],
		})
	}
	return plan
}

type Grid [][]rune

func (g Grid) String() string {
	var s string
	for _, row := range g {
		s += string(row) + "\n"
	}
	return s
}

func (g Grid) Area() int {
	var area int
	for _, row := range g {
		for _, c := range row {
			if c == '#' || c == 'O' {
				area++
			}
		}
	}
	return area
}

type Coord struct {
	r, c int
}

func BuildGrid(plan Plan) (Grid, Coord) {
	// Find the bounds of the grid
	var x, y, x0, x1, y0, y1 int
	for _, instruction := range plan {
		switch instruction.Direction {
		case Up:
			y -= instruction.Length
		case Down:
			y += instruction.Length
		case Left:
			x -= instruction.Length
		case Right:
			x += instruction.Length
		default:
		}
		x0 = min(x0, x)
		x1 = max(x1, x)
		y0 = min(y0, y)
		y1 = max(y1, y)
	}

	// Build the grid
	var grid Grid = make([][]rune, y1-y0+1)
	for i := range grid {
		grid[i] = make([]rune, x1-x0+1)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	// Label the grid
	x, y = -x0, -y0
	for _, instruction := range plan {
		switch instruction.Direction {
		case Up:
			for i := 0; i < instruction.Length; i++ {
				grid[y-i][x] = '#'
			}
			y -= instruction.Length
		case Down:
			for i := 0; i < instruction.Length; i++ {
				grid[y+i][x] = '#'
			}
			y += instruction.Length
		case Left:
			for i := 0; i < instruction.Length; i++ {
				grid[y][x-i] = '#'
			}
			x -= instruction.Length
		case Right:
			for i := 0; i < instruction.Length; i++ {
				grid[y][x+i] = '#'
			}
			x += instruction.Length
		default:
		}
	}
	return grid, Coord{-x0, -y0}
}

func FloodFill(grid Grid, coord Coord) Grid {
	// copy grid
	grid = append([][]rune(nil), grid...)
	for i := range grid {
		grid[i] = append([]rune(nil), grid[i]...)
	}

	q := []Coord{coord}
	for len(q) > 0 {
		coord = q[0]
		q = q[1:]
		if coord.r < 0 || coord.r >= len(grid) || coord.c < 0 || coord.c >= len(grid[coord.r]) {
			continue
		}
		tile := grid[coord.r][coord.c]
		if tile != '#' && (coord.r == 0 || coord.r == len(grid)-1) && (coord.c == 0 || coord.c == len(grid[coord.r])-1) {
			return nil
		}
		if tile != '#' {
			grid[coord.r][coord.c] = '#'
			q = append(q, Coord{coord.r - 1, coord.c})
			q = append(q, Coord{coord.r + 1, coord.c})
			q = append(q, Coord{coord.r, coord.c - 1})
			q = append(q, Coord{coord.r, coord.c + 1})
		}
	}
	return grid
}

func FloodFillEnclosed(grid Grid, origin Coord) Grid {
	if new := FloodFill(grid, Coord{origin.r + 1, origin.c + 1}); new != nil {
		return new
	}
	if new := FloodFill(grid, Coord{origin.r + 1, origin.c - 1}); new != nil {
		return new
	}
	if new := FloodFill(grid, Coord{origin.r - 1, origin.c + 1}); new != nil {
		return new
	}
	if new := FloodFill(grid, Coord{origin.r - 1, origin.c - 1}); new != nil {
		return new
	}
	return nil
}

func part1() {
	lines := ReadLines("input.txt")
	plan := NewPlan(lines)
	grid, origin := BuildGrid(plan)
	grid = FloodFillEnclosed(grid, origin)
	area := grid.Area()
	fmt.Println(area)
}

func part2() {
}

func main() {
	var start time.Time
	start = time.Now()
	part1()
	fmt.Println("Part 1 finished in:", time.Since(start))

	start = time.Now()
	part2()
	fmt.Println("Part 2 finished in:", time.Since(start))
}
