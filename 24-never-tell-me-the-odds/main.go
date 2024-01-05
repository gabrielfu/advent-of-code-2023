package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func ReadFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

func ParseFileLine(fileLine string) (Coord, Coord) {
	var pos, vel Coord
	fmt.Sscanf(
		fileLine,
		"%f, %f, %f @ %f, %f, %f",
		&pos.x, &pos.y, &pos.z, &vel.x, &vel.y, &vel.z,
	)
	return pos, vel
}

type Coord struct {
	x, y, z float64
}

func (c Coord) String() string {
	return fmt.Sprintf("(%.1f, %.1f, %.1f)", c.x, c.y, c.z)
}

type Line2 struct {
	pos Coord
	vel Coord
	m   float64 // slope
	c   float64 // y-intercept
}

func NewLine2(pos Coord, vel Coord) Line2 {
	m := vel.y / vel.x
	c := pos.y - m*pos.x
	return Line2{pos, vel, m, c}
}

func (l Line2) String() string {
	// return fmt.Sprintf("Line(%s)", l.pos)
	return fmt.Sprintf("Line(y = %.2fx + %.2f)", l.m, l.c)
}

func (l Line2) PointIsFuture(point Coord) bool {
	// velocity is always non-zero
	if l.vel.x > 0 {
		return point.x > l.pos.x
	}
	return point.x < l.pos.x
}

func Intersect2(l Line2, o Line2) (Coord, bool) {
	if l.m == o.m {
		return Coord{}, false
	}

	x := (o.c - l.c) / (l.m - o.m)
	y := l.m*x + l.c
	return Coord{x, y, 0}, true
}

func InTestArea2(point Coord, start float64, end float64) bool {
	return point.x >= start && point.x <= end &&
		point.y >= start && point.y <= end
}

func part1() {
	content := ReadFile("input.txt")
	start := 200000000000000.0
	end := 400000000000000.0

	var lines []Line2
	for _, data := range content {
		pos, vel := ParseFileLine(data)
		line := NewLine2(pos, vel)
		lines = append(lines, line)
	}

	total := 0
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			li := lines[i]
			lj := lines[j]
			intersect, ok := Intersect2(li, lj)

			if ok &&
				li.PointIsFuture(intersect) &&
				lj.PointIsFuture(intersect) &&
				InTestArea2(intersect, start, end) {
				total++
			}
		}
	}
	fmt.Println(total)
}

func GaussianElimination(matrix [][]float64) []float64 {
	nr, nc := len(matrix), len(matrix[0])
	r, c := 0, 0
	for r < nr && c < nc {
		max_ := 0.0
		imax := -1
		for i := r; i < nr; i++ {
			cand := math.Abs(matrix[i][c])
			if cand > max_ {
				max_ = cand
				imax = i
			}
		}
		if matrix[imax][c] == 0 {
			c++
		} else {
			matrix[r], matrix[imax] = matrix[imax], matrix[r]
			for i := r + 1; i < nr; i++ {
				factor := matrix[i][c] / matrix[r][c]
				matrix[i][c] = 0
				for j := c + 1; j < nc; j++ {
					matrix[i][j] -= factor * matrix[r][j]
				}
			}
		}
		r++
		c++
	}

	result := make([]float64, nr)
	for i := nr - 1; i >= 0; i-- {
		for j := i + 1; j < nc-1; j++ {
			matrix[i][nc-1] -= matrix[i][j] * result[j]
		}
		result[i] = matrix[i][nc-1] / matrix[i][i]
	}
	return result
}

func CrossProduct(a Coord, b Coord) Coord {
	return Coord{
		a.y*b.z - a.z*b.y,
		a.z*b.x - a.x*b.z,
		a.x*b.y - a.y*b.x,
	}
}

func Diff(a Coord, b Coord) Coord {
	return Coord{
		a.x - b.x,
		a.y - b.y,
		a.z - b.z,
	}
}

func part2() {
	content := ReadFile("input.txt")

	// Px Py Pz Vx Vy Vz t1 t2 t3
	p0, v0 := ParseFileLine(content[0])
	p1, v1 := ParseFileLine(content[1])
	p2, v2 := ParseFileLine(content[2])

	c0 := CrossProduct(p0, v0)
	c1 := CrossProduct(p1, v1)
	c2 := CrossProduct(p2, v2)

	matrix := [][]float64{
		{+(v0.y - v1.y), -(v0.x - v1.x), 0, -(p0.y - p1.y), +(p0.x - p1.x), 0, Diff(c0, c1).z},
		{-(v0.z - v1.z), 0, +(v0.x - v1.x), +(p0.z - p1.z), 0, -(p0.x - p1.x), Diff(c0, c1).y},
		{0, +(v0.z - v1.z), -(v0.y - v1.y), 0, -(p0.z - p1.z), +(p0.y - p1.y), Diff(c0, c1).x},
		{+(v1.y - v2.y), -(v1.x - v2.x), 0, -(p1.y - p2.y), +(p1.x - p2.x), 0, Diff(c1, c2).z},
		{-(v1.z - v2.z), 0, +(v1.x - v2.x), +(p1.z - p2.z), 0, -(p1.x - p2.x), Diff(c1, c2).y},
		{0, +(v1.z - v2.z), -(v1.y - v2.y), 0, -(p1.z - p2.z), +(p1.y - p2.y), Diff(c1, c2).x},
	}

	result := GaussianElimination(matrix)
	fmt.Println(int(math.Round(result[0] + result[1] + result[2])))
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
