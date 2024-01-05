package main

import (
	"fmt"
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
			// if ok {
			// 	fmt.Printf("%s intersects %s at %s\n", li, lj, intersect)
			// } else {
			// 	fmt.Printf("%s does not intersect %s\n", li, lj)
			// }

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
