package main

import (
	"fmt"
	"os"
	"strings"
)

func ReadLines(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

type Node struct {
	Value string
	Left  string
	Right string
}

func (n *Node) String() string {
	return fmt.Sprintf("Node<%s, l=%s, r=%s>", n.Value, n.Left, n.Right)
}

type Graph struct {
	Nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

func (g *Graph) AddNode(value string, left string, right string) {
	g.Nodes[value] = &Node{
		Value: value,
		Left:  left,
		Right: right,
	}
}

func (g *Graph) GetNode(value string) *Node {
	return g.Nodes[value]
}

func part1() {
	lines := ReadLines("input.txt")

	instruction := lines[0]
	g := NewGraph()
	for _, line := range lines[2:] {
		if line == "" {
			continue
		}

		value := line[:3]
		left := line[7:10]
		right := line[12:15]

		g.AddNode(value, left, right)
	}

	cur := "AAA"
	steps := 0
	for {
		direction := instruction[steps%len(instruction)]
		if direction == 'L' {
			cur = g.GetNode(cur).Left
		} else {
			cur = g.GetNode(cur).Right
		}
		steps += 1
		if cur == "ZZZ" {
			break
		}
	}
	println(steps)
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func part2() {
	lines := ReadLines("input.txt")

	instruction := lines[0]
	g := NewGraph()
	for _, line := range lines[2:] {
		if line == "" {
			continue
		}

		value := line[:3]
		left := line[7:10]
		right := line[12:15]

		g.AddNode(value, left, right)
	}

	var nodes []string
	for _, node := range g.Nodes {
		if node.Value[2] == 'A' {
			nodes = append(nodes, node.Value)
		}
	}

	var minSteps []int
	for _, cur := range nodes {
		steps := 0
		for {
			direction := instruction[steps%len(instruction)]
			// fmt.Printf("Node: %s, direction: %c\n", cur, direction)
			if direction == 'L' {
				cur = g.GetNode(cur).Left
			} else {
				cur = g.GetNode(cur).Right
			}
			steps += 1
			if cur[2] == 'Z' {
				break
			}
		}
		minSteps = append(minSteps, steps)
	}

	lcm := LCM(minSteps[0], minSteps[1], minSteps...)
	println(lcm)
}

func main() {
	part1()
	part2()
}
