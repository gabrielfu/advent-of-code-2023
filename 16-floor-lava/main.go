package main

import (
	"fmt"
	"os"
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

type Grid []string

type Coord struct {
	r int
	c int
}

type Direction int

const (
	None Direction = iota
	Up
	Down
	Right
	Left
)

type Beam struct {
	Coord
	Direction
}

func (b *Beam) Next() Beam {
	switch b.Direction {
	case Up:
		return Beam{Coord{b.r - 1, b.c}, b.Direction}
	case Down:
		return Beam{Coord{b.r + 1, b.c}, b.Direction}
	case Right:
		return Beam{Coord{b.r, b.c + 1}, b.Direction}
	case Left:
		return Beam{Coord{b.r, b.c - 1}, b.Direction}
	default:
		return *b
	}
}

func BFS(g Grid, start Beam) int {
	visited := make(map[Beam]struct{})
	visited[start] = struct{}{}

	queue := []Beam{start}

	for len(queue) > 0 {
		var cur Beam
		cur, queue = queue[0], queue[1:]

		var nexts []Beam
		tile := g[cur.r][cur.c]
		switch tile {
		case '.':
			nexts = append(nexts, cur.Next())
		case '|':
			if cur.Direction == Up || cur.Direction == Down {
				nexts = append(nexts, cur.Next())
			} else if cur.Direction == Right || cur.Direction == Left {
				nexts = append(
					nexts,
					Beam{Coord{cur.r - 1, cur.c}, Up},
					Beam{Coord{cur.r + 1, cur.c}, Down},
				)
			}
		case '-':
			if cur.Direction == Right || cur.Direction == Left {
				nexts = append(nexts, cur.Next())
			} else if cur.Direction == Up || cur.Direction == Down {
				nexts = append(
					nexts,
					Beam{Coord{cur.r, cur.c - 1}, Left},
					Beam{Coord{cur.r, cur.c + 1}, Right},
				)
			}
		case '/':
			switch cur.Direction {
			case Up:
				nexts = append(nexts, Beam{Coord{cur.r, cur.c + 1}, Right})
			case Down:
				nexts = append(nexts, Beam{Coord{cur.r, cur.c - 1}, Left})
			case Right:
				nexts = append(nexts, Beam{Coord{cur.r - 1, cur.c}, Up})
			case Left:
				nexts = append(nexts, Beam{Coord{cur.r + 1, cur.c}, Down})
			default:
			}
		case '\\':
			switch cur.Direction {
			case Up:
				nexts = append(nexts, Beam{Coord{cur.r, cur.c - 1}, Left})
			case Down:
				nexts = append(nexts, Beam{Coord{cur.r, cur.c + 1}, Right})
			case Right:
				nexts = append(nexts, Beam{Coord{cur.r + 1, cur.c}, Down})
			case Left:
				nexts = append(nexts, Beam{Coord{cur.r - 1, cur.c}, Up})
			default:
			}
		}

		for _, next := range nexts {
			if next.r < 0 || next.r >= len(g) || next.c < 0 || next.c >= len(g[0]) {
				continue
			}
			if _, ok := visited[next]; ok {
				continue
			}
			visited[next] = struct{}{}
			queue = append(queue, next)
		}
	}

	// get unique coordinates
	energized := make(map[Coord]struct{})
	for b := range visited {
		energized[b.Coord] = struct{}{}
	}

	return len(energized)
}

func part1() {
	grid := ReadLines("input.txt")
	start := Beam{Coord{0, 0}, Right}
	num := BFS(grid, start)
	fmt.Println(num)
}

func part2() {
	grid := ReadLines("input.txt")
	h := len(grid)
	w := len(grid[0])
	maximum := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if r == 0 {
				start := Beam{Coord{r, c}, Down}
				num := BFS(grid, start)
				maximum = max(maximum, num)
			}
			if r == h-1 {
				start := Beam{Coord{r, c}, Up}
				num := BFS(grid, start)
				maximum = max(maximum, num)
			}
			if c == 0 {
				start := Beam{Coord{r, c}, Right}
				num := BFS(grid, start)
				maximum = max(maximum, num)
			}
			if c == w-1 {
				start := Beam{Coord{r, c}, Left}
				num := BFS(grid, start)
				maximum = max(maximum, num)
			}
		}
	}
	println(maximum)
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
