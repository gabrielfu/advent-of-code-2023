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

type Grid struct {
	data []string
	h, w int
}

func NewGrid(lines []string) *Grid {
	g := &Grid{data: lines}
	g.h = len(lines)
	g.w = len(lines[0])
	return g
}

func (g *Grid) String() string {
	return strings.Join(g.data, "\n")
}

func (g *Grid) Get(coord Coord) byte {
	return g.data[coord.r][coord.c]
}

func (g *Grid) Set(coord Coord, ch byte) {
	row := []byte(g.data[coord.r])
	row[coord.c] = ch
	g.data[coord.r] = string(row)
}

func (g *Grid) StartPos() Coord {
	for r, row := range g.data {
		for c, ch := range row {
			if ch == 'S' {
				return Coord{r, c}
			}
		}
	}
	return Coord{-1, -1}
}

type Coord struct {
	r int
	c int
}

type VisitStatus int

const (
	Unvisited VisitStatus = iota
	Odd
	Even
	Unreachable
)

type Entry struct {
	coord Coord
	steps int
}

func traverse(g *Grid, start Coord, steps int) [][]VisitStatus {
	visited := make([][]VisitStatus, g.h)
	for i := range visited {
		visited[i] = make([]VisitStatus, g.w)
	}

	queue := []Entry{{start, 0}}
	var entry Entry
	for len(queue) > 0 {
		entry, queue = queue[0], queue[1:]

		if visited[entry.coord.r][entry.coord.c] != Unvisited {
			continue
		}

		if g.Get(entry.coord) == '#' {
			visited[entry.coord.r][entry.coord.c] = Unreachable
			continue
		}

		switch entry.steps % 2 {
		case 0:
			visited[entry.coord.r][entry.coord.c] = Even
		case 1:
			visited[entry.coord.r][entry.coord.c] = Odd
		}

		if entry.steps == steps {
			continue
		}

		for _, next := range []Coord{
			{entry.coord.r - 1, entry.coord.c},
			{entry.coord.r + 1, entry.coord.c},
			{entry.coord.r, entry.coord.c - 1},
			{entry.coord.r, entry.coord.c + 1},
		} {
			if next.r < 0 || next.r >= g.h || next.c < 0 || next.c >= g.w {
				continue
			}

			if visited[next.r][next.c] != Unvisited {
				continue
			}

			queue = append(queue, Entry{next, entry.steps + 1})
		}
	}
	return visited
}

func Count(visited [][]VisitStatus, status VisitStatus) int {
	total := 0
	for _, row := range visited {
		for _, v := range row {
			if v == status {
				total++
			}
		}
	}
	return total
}

func part1() {
	lines := ReadLines("input.txt")
	grid := NewGrid(lines)
	visited := traverse(grid, grid.StartPos(), 64)
	total := Count(visited, Even)
	fmt.Println(total)
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
