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

func (g Grid) String() string {
	return strings.Join(g, "\n")
}

type Coord struct {
	r int
	c int
}

func findStartingPos(grid Grid) Coord {
	firstRow := grid[0]
	for c, col := range firstRow {
		if col == '.' {
			return Coord{0, c}
		}
	}
	panic("no starting position found")
}

type Node struct {
	Coord
	last  Coord
	steps int
}

func (n Node) String() string {
	return fmt.Sprintf("Node{%v, last=%v, steps=%d}", n.Coord, n.last, n.steps)
}

func solvePart1(grid Grid, start Coord) int {
	h, w := len(grid), len(grid[0])

	var search func(Node) int
	search = func(cur Node) int {
		var neighbors []Coord
		switch grid[cur.r][cur.c] {
		case '>':
			neighbors = []Coord{{cur.r, cur.c + 1}}
		case '<':
			neighbors = []Coord{{cur.r, cur.c - 1}}
		case '^':
			neighbors = []Coord{{cur.r - 1, cur.c}}
		case 'v':
			neighbors = []Coord{{cur.r + 1, cur.c}}
		case '.':
			neighbors = []Coord{
				{cur.r - 1, cur.c},
				{cur.r + 1, cur.c},
				{cur.r, cur.c - 1},
				{cur.r, cur.c + 1},
			}
		}

		best := cur.steps
		for _, nb := range neighbors {
			if nb.r < 0 || nb.r >= h || nb.c < 0 || nb.c >= w {
				continue
			}

			if nb == cur.last {
				continue
			}

			if grid[nb.r][nb.c] == '#' {
				continue
			}

			best = max(best, search(Node{nb, cur.Coord, cur.steps + 1}))
		}
		return best
	}

	return search(Node{start, start, 0})
}

func solvePart2(grid Grid, start Coord) int {
	h, w := len(grid), len(grid[0])

	visited := make(map[Coord]bool)

	var search func(Node) (int, bool)
	search = func(cur Node) (int, bool) {
		if cur.r == h-1 {
			return cur.steps, true
		}

		neighbors := []Coord{
			{cur.r - 1, cur.c},
			{cur.r + 1, cur.c},
			{cur.r, cur.c - 1},
			{cur.r, cur.c + 1},
		}

		best := 0
		ended := false
		for _, nb := range neighbors {
			if nb.r < 0 || nb.r >= h || nb.c < 0 || nb.c >= w {
				continue
			}

			if visited[nb] {
				continue
			}

			if nb == cur.last {
				continue
			}

			if grid[nb.r][nb.c] == '#' {
				continue
			}

			visited[nb] = true
			next, end := search(Node{nb, cur.Coord, cur.steps + 1})
			if end {
				best = max(best, next)
				ended = true
			}
			visited[nb] = false
		}
		return best, ended
	}

	steps, _ := search(Node{start, start, 0})
	return steps
}

func part1() {
	grid := Grid(ReadLines("input.txt"))
	start := findStartingPos(grid)
	steps := solvePart1(grid, start)
	fmt.Println(steps)
}

func part2() {
	// takes 27 minutes to run, probably have a better solution
	grid := Grid(ReadLines("input.txt"))
	start := findStartingPos(grid)
	steps := solvePart2(grid, start)
	fmt.Println(steps)
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
