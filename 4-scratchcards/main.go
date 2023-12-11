package main

import (
	"math"
	"os"
	"strings"
)

func part1() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(content), "\n")

	total := 0
	for _, line := range rows {
		line = strings.Split(line, ":")[1]
		split := strings.Split(line, "|")
		winningNums := strings.Split(strings.TrimSpace(split[0]), " ")
		cardNums := strings.Split(strings.TrimSpace(split[1]), " ")

		intersection := make([]string, 0)
		for _, winningNum := range winningNums {
			if winningNum == "" {
				continue
			}
			for _, cardNum := range cardNums {
				if cardNum == "" {
					continue
				}
				if winningNum == cardNum {
					intersection = append(intersection, winningNum)
				}
			}
		}
		matches := len(intersection)
		if matches > 0 {
			total += int(math.Pow(2, float64(matches-1)))
		}
	}
	println(total)
}

func part2() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(content), "\n")

	copies := make(map[int]int, len(rows))
	for base, line := range rows {
		line = strings.Split(line, ":")[1]
		split := strings.Split(line, "|")
		winningNums := strings.Split(strings.TrimSpace(split[0]), " ")
		cardNums := strings.Split(strings.TrimSpace(split[1]), " ")

		intersection := make([]string, 0)
		for _, winningNum := range winningNums {
			if winningNum == "" {
				continue
			}
			for _, cardNum := range cardNums {
				if cardNum == "" {
					continue
				}
				if winningNum == cardNum {
					intersection = append(intersection, winningNum)
				}
			}
		}
		matches := len(intersection)
		baseCopy := copies[base+1]
		for i := 0; i < matches; i++ {
			copies[base+1+i+1] += baseCopy + 1
		}
	}

	total := len(rows)
	for _, copies := range copies {
		total += copies
	}
	println(total)
}

func main() {
	part1()
	part2()
}
