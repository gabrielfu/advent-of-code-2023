package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func substring(str string, start int, end int) string {
	start = max(start, 0)
	end = min(end, len(str))
	return str[start:end]
}

func part1() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(content), "\n")
	pattern := regexp.MustCompile("[0-9]+")
	total := 0
	for r, line := range rows {
		for _, index := range pattern.FindAllStringIndex(line, -1) {
			left, right := index[0]-1, index[1]
			adjacent := substring(line, left, left+1) +
				substring(line, right, right+1)
			if r > 0 {
				adjacent += substring(rows[r-1], left, right+1)
			}
			if r < len(rows)-1 {
				adjacent += substring(rows[r+1], left, right+1)
			}
			adjacent = strings.ReplaceAll(adjacent, ".", "")
			adjacent = pattern.ReplaceAllString(adjacent, "")
			if len(adjacent) > 0 {
				number, err := strconv.Atoi(line[left+1 : right])
				if err != nil {
					panic(err)
				}
				total += number
			}
		}
	}
	println(total)
}

type Coord struct {
	x int
	y int
}

func part2() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(content), "\n")
	pattern := regexp.MustCompile("[0-9]+")
	coords := make(map[Coord]int)

	for r, line := range rows {
		for _, index := range pattern.FindAllStringIndex(line, -1) {
			left, right := index[0]-1, index[1]
			adjacent := substring(line, left, left+1) +
				substring(line, right, right+1)
			if r > 0 {
				adjacent += substring(rows[r-1], left, right+1)
			}
			if r < len(rows)-1 {
				adjacent += substring(rows[r+1], left, right+1)
			}
			adjacent = strings.ReplaceAll(adjacent, ".", "")
			adjacent = pattern.ReplaceAllString(adjacent, "")
			if len(adjacent) > 0 {
				number, err := strconv.Atoi(line[left+1 : right])
				if err != nil {
					panic(err)
				}
				for x := left + 1; x < right; x++ {
					coords[Coord{x, r}] = number
				}
			}
		}
	}

	ratios := 0
	for r, line := range rows {
		for c, char := range line {
			if char != '*' {
				continue
			}
			adjacent := []Coord{}
			for y := r - 1; y <= r+1; y++ {
				for x := c - 1; x <= c+1; x++ {
					if x == c && y == r {
						continue
					}
					if x < 0 || x >= len(line) || y < 0 || y >= len(rows) {
						continue
					}
					adjacent = append(adjacent, Coord{x, y})
				}
			}
			parts := map[int]struct{}{}
			for _, coord := range adjacent {
				if coords[coord] != 0 {
					parts[coords[coord]] = struct{}{}
				}
			}
			keys := make([]int, 0, len(parts))
			for k := range parts {
				keys = append(keys, k)
			}
			if len(keys) == 2 {
				ratios += keys[0] * keys[1]
			}
		}
	}
	println(ratios)
}

func main() {
	part1()
	part2()
}
