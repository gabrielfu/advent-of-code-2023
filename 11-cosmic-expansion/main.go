package main

import (
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

func (u *Universe) duplicateRow(r int) {
	var data []string
	for i := 0; i < u.Rows(); i++ {
		data = append(data, u.Row(i))
		if i == r {
			data = append(data, u.Row(i))
		}
	}
	u.data = data
}

func (u *Universe) duplicateCol(c int) {
	for i := 0; i < u.Rows(); i++ {
		u.data[i] = u.data[i][:c] + string(u.data[i][c]) + u.data[i][c:]
	}
}

func (u *Universe) Expand() {
	for r := u.Rows() - 1; r >= 0; r-- {
		if !HasGalaxy(u.Row(r)) {
			u.duplicateRow(r)
		}
	}

	for c := u.Cols() - 1; c >= 0; c-- {
		if !HasGalaxy(u.Col(c)) {
			u.duplicateCol(c)
		}
	}
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

func BfsShortestPaths(u *Universe, start Coord) map[Coord]int {
	visited := make(map[Coord]bool)
	queue := []Coord{start}
	steps := 0
	paths := make(map[Coord]int)
	for len(queue) > 0 {
		steps++
		var next []Coord
		for _, coord := range queue {
			if u.IsGalaxy(coord) {
				paths[coord] = steps - 1
			}
			for _, neighbor := range Neighbors(u, coord) {
				if visited[neighbor] {
					continue
				}
				visited[neighbor] = true
				next = append(next, neighbor)
			}
		}
		queue = next
	}
	return paths

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

func part1() {
	lines := ReadLines("input.txt")
	u := NewUniverse(lines)
	u.Expand()

	galaxies := u.Galaxies()
	paths := map[Coord]map[Coord]int{}
	for _, g := range galaxies {
		paths[g] = BfsShortestPaths(u, g)
	}

	pairs := GeneratePairs(galaxies)
	steps := 0
	for _, pair := range pairs {
		steps += paths[pair[0]][pair[1]]
	}
	println(steps)
}

func part2() {
}

func main() {
	part1()
	part2()
}
