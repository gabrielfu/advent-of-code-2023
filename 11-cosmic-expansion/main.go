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

type Universe struct {
	data []string
}

func NewUniverse(data []string) *Universe {
	return &Universe{data: data}
}

func (u *Universe) String() string {
	return strings.Join(u.data, "\n")
}

func (u *Universe) Rows() int {
	return len(u.data)
}

func (u *Universe) Cols() int {
	return len(u.data[0])
}

func (u *Universe) Row(r int) string {
	return u.data[r]
}

func (u *Universe) Col(c int) string {
	var col string
	for _, row := range u.data {
		col += string(row[c])
	}
	return col
}

func HasGalaxy(line string) bool {
	return strings.Contains(line, "#")
}

type Coord struct {
	r, c int
}

func (u *Universe) Galaxies() []Coord {
	var coords []Coord
	for r := 0; r < u.Rows(); r++ {
		for c := 0; c < u.Cols(); c++ {
			if u.data[r][c] == '#' {
				coords = append(coords, Coord{r, c})
			}
		}
	}
	return coords
}

func (u *Universe) EmptyRows() []int {
	var rows []int
	for r := u.Rows() - 1; r >= 0; r-- {
		if !HasGalaxy(u.Row(r)) {
			rows = append(rows, r)
		}
	}
	return rows
}

func (u *Universe) EmptyCols() []int {
	var cols []int
	for c := u.Cols() - 1; c >= 0; c-- {
		if !HasGalaxy(u.Col(c)) {
			cols = append(cols, c)
		}
	}
	return cols
}

type Path struct {
	coord   Coord
	steps   int
	empties int // number of empty rows/cols we have passed
}

func ShortestPaths(u *Universe, start Coord, end Coord, emptyRows, emptyCols []int) Path {
	steps := int(math.Abs(float64(start.r-end.r))) + int(math.Abs(float64(start.c-end.c)))

	// check how many intermediate empty rows/cols we have passed
	empties := 0
	maxR := max(start.r, end.r)
	minR := min(start.r, end.r)
	for _, r := range emptyRows {
		if r > minR && r < maxR {
			empties++
		}
	}
	maxC := max(start.c, end.c)
	minC := min(start.c, end.c)
	for _, c := range emptyCols {
		if c > minC && c < maxC {
			empties++
		}
	}

	return Path{end, steps, empties}

}

func GeneratePairs(coords []Coord) [][]Coord {
	var pairs [][]Coord
	for i, c1 := range coords {
		for _, c2 := range coords[i+1:] {
			pairs = append(pairs, []Coord{c1, c2})
		}
	}
	return pairs
}

func main() {
	lines := ReadLines("input.txt")
	u := NewUniverse(lines)
	emptyRows := u.EmptyRows()
	emptyCols := u.EmptyCols()

	galaxies := u.Galaxies()
	pairs := GeneratePairs(galaxies)
	part1 := 0
	part2 := 0
	for _, pair := range pairs {
		path := ShortestPaths(u, pair[0], pair[1], emptyRows, emptyCols)
		part1 += path.steps + path.empties*(2-1)
		part2 += path.steps + path.empties*(1e6-1)
	}
	println(part1)
	println(part2)
}
