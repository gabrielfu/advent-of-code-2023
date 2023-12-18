package main

import (
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

func AllZero(nums []int) bool {
	for _, num := range nums {
		if num != 0 {
			return false
		}
	}
	return true
}

func Diff(nums []int) []int {
	diffs := []int{}
	for i := 1; i < len(nums); i++ {
		diffs = append(diffs, nums[i]-nums[i-1])
	}
	return diffs
}

func PredictNext(nums []int) int {
	next := 0
	cur := nums
	for {
		cur = Diff(cur)
		if AllZero(cur) {
			break
		}
		next += cur[len(cur)-1]
	}
	return nums[len(nums)-1] + next
}

func PredictPrev(nums []int) int {
	cur := nums
	sub := 0
	mul := -1
	for {
		cur = Diff(cur)
		if AllZero(cur) {
			break
		}
		sub += cur[0] * mul
		mul *= -1
	}
	return nums[0] + sub
}

func part1() {
	lines := ReadLines("input.txt")
	total := 0
	for _, line := range lines {
		parsed := ParseNums(line)
		next := PredictNext(parsed)
		total += next
	}
	println(total)
}

func part2() {
	lines := ReadLines("input.txt")
	total := 0
	for _, line := range lines {
		parsed := ParseNums(line)
		prev := PredictPrev(parsed)
		total += prev
	}
	println(total)
}

func main() {
	part1()
	part2()
}
