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

func (g *Grid) Round(coord Coord) Coord {
	r := coord.r % g.h
	for r < 0 {
		r += g.h
	}
	c := coord.c % g.w
	for c < 0 {
		c += g.w
	}
	return Coord{r, c}
}

func (g *Grid) Get(coord Coord) byte {
	return g.data[coord.r][coord.c]
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

type Entry struct {
	coord Coord
	steps int
}

func traverse(g *Grid, start Coord, steps int) map[int]int {
	res := make(map[int]int)
	visited := make(map[Coord]struct{})
	queue := []Entry{{start, 0}}
	var entry Entry
	for len(queue) > 0 {
		entry, queue = queue[0], queue[1:]

		if entry.steps == steps+1 {
			continue
		}

		if _, ok := visited[entry.coord]; ok {
			continue
		}

		res[entry.steps]++
		visited[entry.coord] = struct{}{}

		for _, next := range []Coord{
			{entry.coord.r - 1, entry.coord.c},
			{entry.coord.r + 1, entry.coord.c},
			{entry.coord.r, entry.coord.c - 1},
			{entry.coord.r, entry.coord.c + 1},
		} {
			if g.Get(g.Round(next)) != '#' {
				queue = append(queue, Entry{next, entry.steps + 1})
			}
		}
	}
	return res
}

func Solve(grid *Grid, n int) int {
	res := traverse(grid, grid.StartPos(), n)
	total := 0
	for dist, num := range res {
		if dist%2 == n%2 {
			total += num
		}
	}
	return total
}

func part1() {
	lines := ReadLines("input.txt")
	grid := NewGrid(lines)
	total := Solve(grid, 64)
	fmt.Println(total)
}

func part2() {
	lines := ReadLines("input.txt")
	grid := NewGrid(lines)

	size := grid.h
	half := size / 2

	y0 := Solve(grid, half)
	y1 := Solve(grid, half+size)
	y2 := Solve(grid, half+size*2)

	// solve polynomial
	a := (y2 - 2*y1 + y0) / 2
	b := (y1 - y0) - a
	c := y0

	target := (26501365 - half) / size
	total := a*target*target + b*target + c
	fmt.Println(total)
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
