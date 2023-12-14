package main

import (
	"fmt"
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

// Find the mirror point P in a sequence of numbers
// such the left and right sides of P are mirror images.
// Returns 0 if no mirror point is found
// E.g., FindMirror([]int{5, 6, 7, 8, 8, 7}) == 4
func FindMirror(arr []int) int {
	for k := 1; k < len(arr); k++ {
		left := arr[:k]
		right := arr[k:]
		// reverse left
		left = append([]int{}, left...)
		for i, j := 0, len(left)-1; i < j; i, j = i+1, j-1 {
			left[i], left[j] = left[j], left[i]
		}
		// fmt.Printf("k=%d, left=%v, right=%v\n", k, left, right)
		end := min(len(left), len(right))
		for i := 0; i < end; i++ {
			if left[i] != right[i] {
				break
			}
			if i == end-1 {
				return k
			}
		}
	}
	return 0
}

func part1() {
	content := ReadFile("input.txt")
	puzzles := strings.Split(strings.TrimSpace(content), "\n\n")

	summary := 0
	for _, puzzle := range puzzles {
		lines := strings.Split(strings.TrimSpace(puzzle), "\n")
		rows := EncodeLines(lines)
		mr := FindMirror(rows)

		transposed := TranposeLines(lines)
		cols := EncodeLines(transposed)
		mc := FindMirror(cols)

		summary += mr*100 + mc
	}
	println(summary)
}

func part2() {
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
