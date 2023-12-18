package main

import (
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type MapRange struct {
	start     int // incl
	end       int // excl
	diff      int
	destStart int
	destEnd   int
}

func MapRangeFromNums(nums []int) MapRange {
	start := nums[1]
	end := nums[1] + nums[2]
	diff := nums[0] - nums[1]
	return MapRange{
		start:     start,
		end:       end,
		diff:      diff,
		destStart: start + diff,
		destEnd:   end + diff,
	}
}

type Map []MapRange

func MapFromNumss(numss [][]int) Map {
	m := Map{}
	for _, nums := range numss {
		m = append(m, MapRangeFromNums(nums))
	}
	Sort(m)
	return m
}

func (m Map) Translate(num int) int {
	for _, mr := range m {
		if num >= mr.start && num < mr.end {
			return num + mr.diff
		}
	}
	return num
}

func (m Map) BackTranslate(num int) int {
	for _, mr := range m {
		if num >= mr.destStart && num < mr.destEnd {
			return num - mr.diff
		}
	}
	return num
}

func Sort[Map ~[]MapRange](m Map) {
	sort.Slice(m, func(i, j int) bool {
		return m[i].start < m[j].start
	})
}

type Maps []Map

func (m Maps) Translate(num int) int {
	for _, m := range m {
		num = m.Translate(num)
	}
	return num
}

func (m Maps) BackTranslate(num int) int {
	for i := len(m) - 1; i >= 0; i-- {
		num = m[i].BackTranslate(num)
	}
	return num
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

type SeedRange struct {
	start int
	end   int
}

func part1() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(content), "\n")

	// parse seeds
	seeds := ParseNums(strings.Split(rows[0], ":")[1])

	// build maps
	maps := Maps{}
	numss := [][]int{}
	for _, row := range rows[1:] {
		if row == "" || !unicode.IsDigit(rune(row[0])) {
			if len(numss) > 0 {
				maps = append(maps, MapFromNumss(numss))
			}
			numss = [][]int{}
			continue
		}
		nums := ParseNums(row)
		if len(nums) > 0 {
			numss = append(numss, nums)
		}
	}
	maps = append(maps, MapFromNumss(numss))

	// translate seeds
	locations := make([]int, len(seeds))
	for i, seed := range seeds {
		locations[i] = maps.Translate(seed)
	}

	println(slices.Min(locations))
}

func part2() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(content), "\n")

	// parse seeds
	seeds := ParseNums(strings.Split(rows[0], ":")[1])
	ranges := []SeedRange{}
	for i := 0; i < len(seeds); i += 2 {
		ranges = append(ranges, SeedRange{seeds[i], seeds[i] + seeds[i+1]})
	}

	// build maps
	maps := Maps{}
	numss := [][]int{}
	for _, row := range rows[1:] {
		if row == "" || !unicode.IsDigit(rune(row[0])) {
			if len(numss) > 0 {
				maps = append(maps, MapFromNumss(numss))
			}
			numss = [][]int{}
			continue
		}
		nums := ParseNums(row)
		if len(nums) > 0 {
			numss = append(numss, nums)
		}
	}
	maps = append(maps, MapFromNumss(numss))

	for location := 0; ; location++ {
		seed := maps.BackTranslate(location)
		for _, seedRange := range ranges {
			if seed >= seedRange.start && seed < seedRange.end {
				println(location)
				return
			}
		}
	}
}

func main() {
	part1()
	part2()
}
