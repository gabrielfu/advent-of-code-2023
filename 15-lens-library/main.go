package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ReadFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func HashChar(char rune, cur int) int {
	cur += int(char)
	cur *= 17
	cur %= 256
	return cur
}

func HashString(s string) int {
	cur := 0
	for _, char := range s {
		cur = HashChar(char, cur)
	}
	return cur
}

func part1() {
	line := ReadFile("input.txt")
	strs := strings.Split(line, ",")

	total := 0
	for _, str := range strs {
		total += HashString(str)
	}
	println(total)
}

type Lens struct {
	label       string
	hash        int
	focalLength int
	i           int
}

func part2() {
	line := ReadFile("input.txt")
	strs := strings.Split(line, ",")
	pattern := regexp.MustCompile(`([a-z]+)(=|-)([0-9]*)`)

	boxes := make([]map[string]*Lens, 256)
	for i, str := range strs {
		matches := pattern.FindAllStringSubmatch(str, -1)
		label := matches[0][1]
		b := HashString(label)
		op := matches[0][2]

		if boxes[b] == nil {
			boxes[b] = make(map[string]*Lens)
		}

		switch op {
		case "=":
			fl, err := strconv.Atoi(matches[0][3])
			if err != nil {
				panic(err)
			}
			if boxes[b][label] == nil {
				boxes[b][label] = &Lens{label, b, fl, i}
			} else {
				boxes[b][label].focalLength = fl
			}
		case "-":
			delete(boxes[b], label)
		default:
		}
	}

	total := 0
	for b, box := range boxes {
		lenses := []*Lens{}
		for _, lens := range box {
			lenses = append(lenses, lens)
		}
		sort.Slice(lenses, func(i, j int) bool {
			return lenses[i].i < lenses[j].i
		})
		for i, lens := range lenses {
			total += (b + 1) * (i + 1) * lens.focalLength
		}
	}
	println(total)
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
