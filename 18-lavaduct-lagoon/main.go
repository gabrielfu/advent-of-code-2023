package main

import (
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

type Direction string

const (
	Up    Direction = "U"
	Down  Direction = "D"
	Left  Direction = "L"
	Right Direction = "R"
)

type PlanItem struct {
	Direction Direction
	Length    int
	Color     string
}

func (i PlanItem) Decode() PlanItem {
	hex := i.Color[2:7]
	length, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		panic(err)
	}

	d := i.Color[7]
	var direction Direction
	switch d {
	case '0':
		direction = Right
	case '1':
		direction = Down
	case '2':
		direction = Left
	case '3':
		direction = Up
	default:
	}

	return PlanItem{
		Direction: direction,
		Length:    int(length),
		Color:     i.Color,
	}
}

type Plan []PlanItem

func (p Plan) String() string {
	var s string
	for _, inst := range p {
		s += fmt.Sprintf("%s %d %s\n", inst.Direction, inst.Length, inst.Color)
	}
	return s
}

func (p Plan) TotalLength() int {
	var t int
	for _, inst := range p {
		t += inst.Length
	}
	return t
}

func NewPlan(lines []string, decode bool) Plan {
	var plan Plan
	for _, line := range lines {
		split := strings.Split(line, " ")
		l, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}
		inst := PlanItem{
			Direction: Direction(split[0]),
			Length:    l,
			Color:     split[2],
		}
		if decode {
			inst = inst.Decode()
		}
		plan = append(plan, inst)
	}
	return plan
}

type Coord struct {
	r, c int
}

func GetCoordList(plan Plan) []Coord {
	coord := Coord{0, 0}
	result := []Coord{coord}
	for _, instruction := range plan {
		switch instruction.Direction {
		case Up:
			coord = Coord{coord.r - instruction.Length, coord.c}
		case Down:
			coord = Coord{coord.r + instruction.Length, coord.c}
		case Left:
			coord = Coord{coord.r, coord.c - instruction.Length}
		case Right:
			coord = Coord{coord.r, coord.c + instruction.Length}
		default:
		}
		result = append(result, coord)
	}
	return result
}

// Calculates the area using Shoelace formula
func ShoelaceArea(coords []Coord) int {
	var area int
	for i := 0; i < len(coords)-1; i++ {
		this := coords[i]
		next := coords[i+1]
		area += (this.r + next.r) * (this.c - next.c)
	}
	return area / 2
}

func part1() {
	lines := ReadLines("input.txt")
	plan := NewPlan(lines, false)
	coords := GetCoordList(plan)
	area := ShoelaceArea(coords)
	length := plan.TotalLength()
	totalArea := area + length/2 + 1
	fmt.Println(totalArea)
}

func part2() {
	lines := ReadLines("input.txt")
	plan := NewPlan(lines, true)
	coords := GetCoordList(plan)
	area := ShoelaceArea(coords)
	length := plan.TotalLength()
	totalArea := area + length/2 + 1
	fmt.Println(totalArea)
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
