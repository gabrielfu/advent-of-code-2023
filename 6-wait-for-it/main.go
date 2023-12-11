package main

import (
	"math"
	"os"
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

func ParseNums(line string) []int {
	nums := []int{}
	for _, num := range strings.Split(line, " ") {
		if num == "" {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			panic(err)
		}
		nums = append(nums, n)
	}
	return nums
}

func part1() {
	lines := ReadLines("input.txt")
	times := ParseNums(strings.Split(lines[0], ":")[1])
	distances := ParseNums(strings.Split(lines[1], ":")[1])

	// Given time is T, distance is D, record is R
	// and we hold the button for n milliseconds, we have the constraint:
	// D = (T - n) * n > R
	// So the root for n is
	// n = (T ± sqrt(T^2 - 4R)) / 2
	ways := 1
	for i, t := range times {
		d := distances[i]
		discriminant := float64(t*t - 4*d)
		small := math.Ceil((float64(t) - math.Sqrt(discriminant)) / 2)
		large := math.Floor((float64(t) + math.Sqrt(discriminant)) / 2)
		ways *= int(large) - int(small) + 1
	}
	println(ways)
}

func part2() {
	lines := ReadLines("input.txt")

	t, err := strconv.Atoi(strings.ReplaceAll(strings.Split(lines[0], ":")[1], " ", ""))
	if err != nil {
		panic(err)
	}
	d, err := strconv.Atoi(strings.ReplaceAll(strings.Split(lines[1], ":")[1], " ", ""))
	if err != nil {
		panic(err)
	}

	discriminant := float64(t*t - 4*d)
	small := math.Ceil((float64(t) - math.Sqrt(discriminant)) / 2)
	large := math.Floor((float64(t) + math.Sqrt(discriminant)) / 2)
	ways := int(large) - int(small) + 1
	println(ways)
}

func main() {
	part1()
	part2()
}
