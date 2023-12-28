package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
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

type Coord struct {
	x, y, z int
}

type Id int

type Brick struct {
	Start Coord
	End   Coord
	Id    Id
}

func (b *Brick) String() string {
	return fmt.Sprintf(
		"Brick(%d: %d,%d,%d~%d,%d,%d)",
		b.Id, b.Start.x, b.Start.y, b.Start.z, b.End.x, b.End.y, b.End.z,
	)
}

func (b *Brick) MoveDown() *Brick {
	return &Brick{
		Coord{b.Start.x, b.Start.y, b.Start.z - 1},
		Coord{b.End.x, b.End.y, b.End.z - 1},
		b.Id,
	}
}

func (b *Brick) Intersects(other *Brick) bool {
	return b.Start.x <= other.End.x && b.End.x >= other.Start.x &&
		b.Start.y <= other.End.y && b.End.y >= other.Start.y &&
		b.Start.z <= other.End.z && b.End.z >= other.Start.z
}

func ParseBrick(line string, id int) *Brick {
	var start, end Coord
	fmt.Sscanf(
		line,
		"%d,%d,%d~%d,%d,%d",
		&start.x, &start.y, &start.z, &end.x, &end.y, &end.z,
	)
	return &Brick{start, end, Id(id)}
}

func ParseBricks(lines []string) []*Brick {
	bricks := make([]*Brick, len(lines))
	for i, line := range lines {
		bricks[i] = ParseBrick(line, i)
	}
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].Start.z < bricks[j].Start.z
	})
	return bricks
}

func HasIntersecting(b *Brick, others []*Brick) bool {
	for _, o := range others {
		if b.Intersects(o) {
			return true
		}
	}
	return false
}

func GetBricks(bricks []*Brick, ids []Id) []*Brick {
	var res []*Brick
	for _, b := range bricks {
		if slices.Contains(ids, b.Id) {
			res = append(res, b)
		}
	}
	return res
}

func StartFalling(bricks []*Brick) {
	// map of z-level to brick ids
	zmap := make(map[int][]Id)
	for i := range bricks {
		for {
			brick := bricks[i]
			cand := brick.MoveDown()
			z := cand.Start.z

			if z < 1 || HasIntersecting(cand, GetBricks(bricks, zmap[z])) {
				break
			}

			bricks[i] = cand
		}

		brick := bricks[i]
		for z := brick.Start.z; z <= brick.End.z; z++ {
			zmap[z] = append(zmap[z], brick.Id)
		}
	}
}

// key is supported by values
func BuildSupportedGraph(bricks []*Brick) map[Id][]Id {
	graph := make(map[Id][]Id)
	for i, b1 := range bricks {
		for j, b2 := range bricks {
			if i == j {
				continue
			}
			down := &Brick{
				Coord{b1.Start.x, b1.Start.y, b1.Start.z - 1},
				Coord{b1.End.x, b1.End.y, b1.End.z - 1},
				b1.Id,
			}
			if down.Intersects(b2) {
				graph[b1.Id] = append(graph[b1.Id], b2.Id)
			}
		}
	}
	return graph
}

func InverseGraph(graph map[Id][]Id) map[Id][]Id {
	inv := make(map[Id][]Id)
	for k, vs := range graph {
		for _, v := range vs {
			inv[v] = append(inv[v], k)
		}
	}
	return inv
}

// key is supporting values
func BuildSupportingGraph(bricks []*Brick) map[Id][]Id {
	g := BuildSupportedGraph(bricks)
	return InverseGraph(g)
}

func FindDisintegrable(bricks []*Brick) []Id {
	graph := BuildSupportedGraph(bricks)

	// find all bricks that are not supported by any other brick
	unsafe := map[Id]struct{}{}
	for _, supported := range graph {
		if len(supported) == 1 {
			unsafe[supported[0]] = struct{}{}
		}
	}

	var disintegrable []Id
	for _, bricks := range bricks {
		if _, ok := unsafe[bricks.Id]; !ok {
			disintegrable = append(disintegrable, bricks.Id)
		}
	}
	return disintegrable
}

func containsAll[T comparable](set map[T]struct{}, subset []T) bool {
	for _, v := range subset {
		if _, ok := set[v]; !ok {
			return false
		}
	}
	return true
}

func addAll[T comparable](set map[T]struct{}, elems []T) {
	for _, v := range elems {
		set[v] = struct{}{}
	}
}

func removeBrick(supported, supporting map[Id][]Id, id Id, removed map[Id]struct{}) {
	if supporting[id] == nil || len(supporting[id]) == 0 {
		return
	}

	removed[id] = struct{}{}
	var nexts []Id
	for _, supp := range supporting[id] {
		if containsAll(removed, supported[supp]) {
			nexts = append(nexts, supp)
		}
	}

	addAll(removed, nexts)
	for _, next := range nexts {
		removeBrick(supported, supporting, next, removed)
	}
}

func RemoveBrick(supported, supporting map[Id][]Id, id Id) int {
	removed := make(map[Id]struct{})
	removeBrick(supported, supporting, id, removed)
	delete(removed, id)
	return len(removed)
}

func part1() {
	lines := ReadLines("input.txt")
	bricks := ParseBricks(lines)
	StartFalling(bricks)
	disintegrable := FindDisintegrable(bricks)
	fmt.Println(len(disintegrable))
}

func part2() {
	lines := ReadLines("input.txt")
	bricks := ParseBricks(lines)
	StartFalling(bricks)

	supported := BuildSupportedGraph(bricks)
	supporting := BuildSupportingGraph(bricks)
	total := 0
	for _, brick := range bricks {
		total += RemoveBrick(supported, supporting, brick.Id)
	}
	fmt.Println(total)
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
