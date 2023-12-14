package main

import (
	"fmt"
	"math/bits"
	"os"
	"strings"
	"time"
)

func ReadFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(content)
}

func TranposeLines(lines []string) []string {
	var transposed []string
	for i := 0; i < len(lines[0]); i++ {
		var row []byte
		for _, line := range lines {
			row = append(row, line[i])
		}
		transposed = append(transposed, string(row))
	}
	return transposed
}

// Encode a string as a binary number, where '#' is 1 and '.' is 0.
func Encode(s string) int {
	h := 0
	for _, char := range s {
		h <<= 1
		if char == '#' {
			h |= 1
		}
	}
	return h
}

func EncodeLines(lines []string) []int {
	var encoded []int
	for _, line := range lines {
		encoded = append(encoded, Encode(line))
	}
	return encoded
}

func ReverseSlice(arr []int) []int {
	arr = append([]int{}, arr...)
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func NumSmudges(a, b int) int {
	return bits.OnesCount(uint(a ^ b))
}

func FindMirror(arr []int, smudge bool) (int, bool) {
	for k := 1; k < len(arr); k++ {
		count := 0
		left := arr[:k]
		right := arr[k:]
		left = ReverseSlice(left)
		end := min(len(left), len(right))
		for i := 0; i < end; i++ {
			count += NumSmudges(left[i], right[i])
			if i == end-1 && ((smudge && count == 1) || (!smudge && count == 0)) {
				return k, true
			}
		}
	}
	return 0, false
}

func Solve(smudge bool) {
	content := ReadFile("input.txt")
	puzzles := strings.Split(strings.TrimSpace(content), "\n\n")

	summary := 0
	for _, puzzle := range puzzles {
		lines := strings.Split(strings.TrimSpace(puzzle), "\n")
		rows := EncodeLines(lines)
		mr, ok := FindMirror(rows, smudge)
		if ok {
			summary += mr * 100
		}

		transposed := TranposeLines(lines)
		cols := EncodeLines(transposed)
		mc, ok := FindMirror(cols, smudge)
		if ok {
			summary += mc
		}
	}
	println(summary)
}

func part1() {
	Solve(false)
}

func part2() {
	Solve(true)
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
