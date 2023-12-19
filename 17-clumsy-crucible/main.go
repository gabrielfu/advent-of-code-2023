package main

import (
	"container/heap"
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

func ParseNums(line string) []int {
	nums := []int{}
	for _, num := range strings.Split(line, "") {
		if num == "" {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			panic(err)
		}
		nums = append(nums, n)
	}
	return nums
}

func BuildGrid(lines []string) [][]int {
	grid := make([][]int, len(lines))
	for i, line := range lines {
		grid[i] = ParseNums(line)
	}
	return grid
}

type Coord struct {
	r, c int
}

func (c *Coord) Move(dir Direction) *Coord {
	switch dir {
	case Up:
		return &Coord{c.r - 1, c.c}
	case Down:
		return &Coord{c.r + 1, c.c}
	case Left:
		return &Coord{c.r, c.c - 1}
	case Right:
		return &Coord{c.r, c.c + 1}
	default:
		return nil
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
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Right:
		return "Right"
	case Left:
		return "Left"
	default:
		return "None"
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

type Node struct {
	pos   Coord
	dir   Direction
	steps int
}

type Crucible struct {
	pos   Coord
	dir   Direction
	cost  int
	steps int // consecutive steps in the same direction
}

func (c *Crucible) String() string {
	return fmt.Sprintf("Crucible(pos=%v dir=%s cost=%v steps=%d)", c.pos, c.dir, c.cost, c.steps)
}

func (c *Crucible) Node() Node {
	return Node{c.pos, c.dir, c.steps}
}

type PriorityQueue []*Crucible

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Crucible))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	x := old[0]
	*pq = old[1:]
	return x
}

func Neighbors(grid [][]int, c *Crucible, minSteps, maxSteps int) []*Crucible {
	neighbors := []*Crucible{}
	h, w := len(grid), len(grid[0])
	for _, dir := range []Direction{Up, Down, Left, Right} {
		if c.dir == dir && c.steps >= maxSteps {
			continue
		}
		if c.dir != dir && c.steps > 0 && c.steps < minSteps {
			continue
		}
		if c.dir == dir.Opposite() {
			continue
		}

		next := c.pos.Move(dir)
		if next == nil {
			continue
		}
		if next.r < 0 || next.r >= h || next.c < 0 || next.c >= w {
			continue
		}
		alt := c.cost + grid[next.r][next.c]
		steps := 1
		if c.dir == dir {
			steps = c.steps + 1
		}
		neighbors = append(neighbors, &Crucible{*next, dir, alt, steps})
	}
	return neighbors
}

func Dijkstra(grid [][]int, source, dest Coord, minSteps, maxSteps int) int {
	pq := PriorityQueue{&Crucible{source, None, 0, 0}}
	visited := map[Node]struct{}{}

	minimum := int(1e9)
	for len(pq) > 0 {
		cur := pq.Pop().(*Crucible)

		if cur.pos == dest && cur.steps >= minSteps {
			minimum = min(minimum, cur.cost)
		}

		for _, nb := range Neighbors(grid, cur, minSteps, maxSteps) {
			node := nb.Node()
			if _, ok := visited[node]; !ok {
				visited[node] = struct{}{}
				pq.Push(nb)
			}
		}
		heap.Init(&pq)
	}
	return minimum
}

func part1() {
	lines := ReadLines("input.txt")
	grid := BuildGrid(lines)
	h, w := len(grid), len(grid[0])
	cost := Dijkstra(grid, Coord{0, 0}, Coord{h - 1, w - 1}, 1, 3)
	fmt.Println(cost)
}

func part2() {
	lines := ReadLines("input.txt")
	grid := BuildGrid(lines)
	h, w := len(grid), len(grid[0])
	cost := Dijkstra(grid, Coord{0, 0}, Coord{h - 1, w - 1}, 4, 10)
	fmt.Println(cost)
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
