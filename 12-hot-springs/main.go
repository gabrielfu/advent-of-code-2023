package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ReadLines(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

var AllConditions = []string{".", "#"}

func ArragementGenerator(conditions string, unknown []int) []string {
	if len(unknown) == 0 {
		return []string{conditions}
	}
	index := unknown[0]
	unknown = unknown[1:]

	var result []string
	for _, c := range AllConditions {
		new := conditions[:index] + c + conditions[index+1:]
		result = append(result, ArragementGenerator(new, unknown)...)
	}
	return result
}

func UnknownIndex(conditions string) []int {
	var result []int
	for i, c := range conditions {
		if c == '?' {
			result = append(result, i)
		}
	}
	return result
}

var DamagedPattern = regexp.MustCompile(`#+`)

func ProduceRecord(conditions string) string {
	matches := DamagedPattern.FindAllStringIndex(conditions, -1)
	var groups []string
	for _, m := range matches {
		groups = append(groups, strconv.Itoa(m[1]-m[0]))
	}
	return strings.Join(groups, ",")
}

func part1() {
	lines := ReadLines("input.txt")
	ways := 0
	for _, line := range lines {
		split := strings.Split(line, " ")
		conditions := split[0]
		record := split[1]

		unknown := UnknownIndex(conditions)
		arrangements := ArragementGenerator(conditions, unknown)
		for _, a := range arrangements {
			if ProduceRecord(a) == record {
				ways++
			}
		}
	}
	println(ways)
}

func part2() {
}

func main() {
	part1()
	part2()
}
