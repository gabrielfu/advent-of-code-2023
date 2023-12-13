package main

import (
	"fmt"
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

func (u *Universe) IsGalaxy(coord Coord) bool {
	return u.data[coord.r][coord.c] == '#'
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

var dr = []int{-1, 1, 0, 0}
var dc = []int{0, 0, -1, 1}

func Neighbors(u *Universe, coord Coord) []Coord {
	var neighbors []Coord
	for i := 0; i < 4; i++ {
		r, c := dr[i], dc[i]
		new := Coord{coord.r + r, coord.c + c}
		if new.r < 0 || new.r >= u.Rows() || new.c < 0 || new.c >= u.Cols() {
			continue
		}
		neighbors = append(neighbors, new)
	}
	return neighbors
}

type Path struct {
	coord   Coord
	steps   int
	empties int // number of empty rows/cols we have passed
}

func BfsShortestPaths(u *Universe, start Coord, emptyRows, emptyCols []int) map[Coord]Path {
	visited := make(map[Coord]bool)
	queue := []Path{{start, 0, 0}}
	result := make(map[Coord]Path)
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		visited[path.coord] = true
		for _, neighbor := range Neighbors(u, path.coord) {
			if visited[neighbor] {
				continue
			}
			if u.IsGalaxy(neighbor) {
				result[neighbor] = Path{neighbor, path.steps + 1, path.empties}
			}
			empties := path.empties
			if slices.Contains(emptyRows, neighbor.r) || slices.Contains(emptyCols, neighbor.c) {
				empties++
			}
			queue = append(queue, Path{neighbor, path.steps + 1, empties})
		}
	}
	return result

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
	paths := map[Coord]map[Coord]Path{}
	for i, g := range galaxies {
		paths[g] = BfsShortestPaths(u, g, emptyRows, emptyCols)
		fmt.Printf("Completed %d/%d\n", i+1, len(galaxies))
	}

	pairs := GeneratePairs(galaxies)
	part1 := 0
	part2 := 0
	for _, pair := range pairs {
		path := paths[pair[0]][pair[1]]
		part1 += path.steps + path.empties*(2-1)
		part2 += path.steps + path.empties*(1e6-1)
	}
	println(part1)
	println(part2)
}
