package main

import (
	"os"
	"strconv"
	"strings"
)

func part1() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	availableCubes := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	sum := 0
	for _, line := range lines {
		splits := strings.Split(line, ":")
		prefix, game := splits[0], splits[1]
		id, err := strconv.Atoi(strings.Split(prefix, " ")[1])
		if err != nil {
			panic(err)
		}

		sets := strings.Split(game, ";")
		ok := true
		for _, set := range sets {
			if !ok {
				break
			}
			cubes := strings.Split(set, ",")
			for _, cube := range cubes {
				splits = strings.Split(strings.TrimSpace(cube), " ")
				num, err := strconv.Atoi(strings.TrimSpace(splits[0]))
				if err != nil {
					panic(err)
				}
				color := strings.TrimSpace(splits[1])
				if num > availableCubes[color] {
					ok = false
					break
				}
			}
		}

		if ok {
			sum += id
		}
	}

	println(sum)
}

func part2() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")

	powers := 0
	for _, line := range lines {
		splits := strings.Split(line, ":")
		game := splits[1]
		sets := strings.Split(game, ";")
		cubeMap := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}
		for _, set := range sets {
			cubes := strings.Split(set, ",")
			for _, cube := range cubes {
				splits = strings.Split(strings.TrimSpace(cube), " ")
				num, err := strconv.Atoi(strings.TrimSpace(splits[0]))
				if err != nil {
					panic(err)
				}
				color := strings.TrimSpace(splits[1])
				cubeMap[color] = max(cubeMap[color], num)
			}
		}

		power := 1
		for _, value := range cubeMap {
			power *= value
		}

		powers += power
	}

	println(powers)
}

func main() {
	part1()
	part2()
}
