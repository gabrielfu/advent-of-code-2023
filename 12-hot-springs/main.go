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

func SplitNums(line string) []int {
	var result []int
	for _, s := range strings.Split(line, ",") {
		n, _ := strconv.Atoi(s)
		result = append(result, n)
	}
	return result
}

var cache = make(map[string]int)

func cacheKey(conditions string, groups []int) string {
	key := conditions
	for _, g := range groups {
		key += strconv.Itoa(g) + ","
	}
	return key
}

func CountArrangements(conditions string, groups []int) int {
	key := cacheKey(conditions, groups)
	if v, ok := cache[key]; ok {
		return v
	}

	// base cases
	if len(conditions) == 0 {
		if len(groups) == 0 {
			return 1
		} else {
			return 0
		}
	}

	switch conditions[0] {
	case '?':
		good := CountArrangements("."+conditions[1:], groups)
		damaged := CountArrangements("#"+conditions[1:], groups)
		return good + damaged
	case '.':
		count := CountArrangements(conditions[1:], groups)
		cache[key] = count
		return count
	case '#':
		if len(groups) == 0 {
			cache[key] = 0
			return 0
		}
		size := groups[0]
		if len(conditions) < size {
			cache[key] = 0
			return 0
		}
		if strings.Contains(conditions[:size], ".") {
			cache[key] = 0
			return 0
		}

		if len(groups) == 1 {
			count := CountArrangements(conditions[size:], []int{})
			cache[key] = count
			return count
		}

		if len(conditions) < size+1 || conditions[size] == '#' {
			cache[key] = 0
			return 0
		}

		count := CountArrangements(conditions[size+1:], groups[1:])
		cache[key] = count
		return count
	default:
		return 0
	}
}

func part1() {
	lines := ReadLines("input.txt")
	ways := 0
	for _, line := range lines {
		split := strings.Split(line, " ")
		conditions := split[0]
		groups := SplitNums(split[1])
		counts := CountArrangements(conditions, groups)
		ways += counts
	}
	println(ways)
}

func repeatSlice[T any](s []T, n int) []T {
	arr := make([]T, len(s)*n)
	for i := 0; i < n; i++ {
		copy(arr[i*len(s):], s)
	}
	return arr
}

func part2() {
	lines := ReadLines("input.txt")
	ways := 0
	for _, line := range lines {
		split := strings.Split(line, " ")
		conditions := split[0]
		groups := SplitNums(split[1])

		// unfold 5 times and join by "?"
		conditions = strings.Repeat("?"+conditions, 5)[1:]
		groups = repeatSlice(groups, 5)

		counts := CountArrangements(conditions, groups)
		ways += counts
	}
	println(ways)
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
