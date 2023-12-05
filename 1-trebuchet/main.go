package main

import (
	"os"
	"regexp"
	"strings"
)

func part1() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	r := regexp.MustCompile("[^0-9]+")

	sum := 0
	for _, line := range lines {
		digits := r.ReplaceAllString(line, "")
		if len(digits) == 0 {
			panic("line has no digits: " + line)
		}
		tens, ones := digits[0], digits[len(digits)-1]
		sum += 10*int(tens-'0') + int(ones-'0')
	}
	println(sum)
}

func part2() {
	digitMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
		"1":     1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
	}

	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	sum := 0
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		tens, ones := 0, 0
		tensIndex, onesIndex := len(line), -1
		for k, v := range digitMap {
			i := strings.Index(line, k)
			if i != -1 && i < tensIndex {
				tensIndex = i
				tens = v
			}
			j := strings.LastIndex(line, k)
			if j != -1 && j > onesIndex {
				onesIndex = j
				ones = v
			}
		}
		sum += 10*tens + ones
	}
	println(sum)
}

func main() {
	part1()
	part2()
}
