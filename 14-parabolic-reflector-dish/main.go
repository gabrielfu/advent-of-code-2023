package main

import (
	"fmt"
	"os"
	"slices"
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

func FormatLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func RotateClockwise(lines []string) []string {
	var rotated []string
	for i := 0; i < len(lines[0]); i++ {
		var row []byte
		for j := len(lines) - 1; j >= 0; j-- {
			row = append(row, lines[j][i])
		}
		rotated = append(rotated, string(row))
	}
	return rotated
}

func RotateAnticlockwise(lines []string) []string {
	var rotated []string
	for i := len(lines[0]) - 1; i >= 0; i-- {
		var row []byte
		for j := 0; j < len(lines); j++ {
			row = append(row, lines[j][i])
		}
		rotated = append(rotated, string(row))
	}
	return rotated
}

// Roll the round rocks to the right for one line
func RollLine(line string) string {
	N := len(line)
	config := make([]int, N)
	rounds := 0
	for i, char := range line {
		switch char {
		case 'O':
			// pick up the rock
			rounds++
		case '#':
			config[i] = -1
			if i >= 1 && line[i-1] != '#' {
				config[i-1] = rounds
			}
			rounds = 0
		default:
		}
		if i == N-1 && line[i] != '#' {
			config[i] = rounds
		}
	}

	rolled := make([]byte, len(line))
	for i := N - 1; i >= 0; i-- {
		if config[i] > 0 {
			rolled[i] = 'O'
			// drop one rock. config[i-1] should never be -1 (#) if the input is correct
			if i >= 1 && config[i] > 1 {
				config[i-1] = config[i] - 1
				config[i] = 1
			}
		} else if config[i] == -1 {
			rolled[i] = '#'
		} else {
			rolled[i] = '.'
		}
	}
	return string(rolled)
}

func RollPlatform(lines []string) []string {
	rolled := make([]string, len(lines))
	for i, line := range lines {
		rolled[i] = RollLine(line)
	}
	return rolled
}

func RollOneCycle(lines []string) []string {
	var rotated []string
	rotated = append(rotated, lines...)
	for range [4]int{} {
		rotated = RotateClockwise(rotated)
		rotated = RollPlatform(rotated)
	}
	return rotated
}

func ScoreRow(line string) int {
	score := 0
	for i, char := range line {
		if char == 'O' {
			score += i + 1
		}
	}
	return score
}

func ScorePlatform(lines []string) int {
	score := 0
	lines = RotateClockwise(lines)
	for _, line := range lines {
		score += ScoreRow(line)
	}
	return score
}

func part1() {
	lines := ReadLines("input.txt")
	lines = RotateClockwise(lines)
	lines = RollPlatform(lines)
	lines = RotateAnticlockwise(lines)
	score := ScorePlatform(lines)
	println(score)
}

func part2() {
	lines := ReadLines("input.txt")

	var cycles []string
	var scores []int
	n := 1000000000
	loopStart := -1
	i := 0
	for ; i < n; i++ {
		lines = RollOneCycle(lines)
		formatted := FormatLines(lines)

		loopStart = slices.Index(cycles, formatted)
		if loopStart != -1 {
			break
		}
		cycles = append(cycles, formatted)
		scores = append(scores, ScorePlatform(lines))
	}
	index := (n-loopStart)%(i-loopStart+1) + loopStart
	println(scores[index])
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
