package main

import (
	"fmt"
	"math/rand"
	"os"
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

type Node string

type Graph struct {
	nodes []Node
	adj   map[Node][]Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes: []Node{},
		adj:   make(map[Node][]Node),
	}
}

func (g *Graph) RandomNode() Node {
	return g.nodes[rand.Intn(len(g.nodes))]
}

func (g *Graph) HasNode(n Node) bool {
	for _, node := range g.nodes {
		if node == n {
			return true
		}
	}
	return false
}

func (g *Graph) AddEdge(a, b Node) {
	if !g.HasNode(a) {
		g.nodes = append(g.nodes, a)
	}
	if !g.HasNode(b) {
		g.nodes = append(g.nodes, b)
	}
	g.adj[a] = append(g.adj[a], b)
	g.adj[b] = append(g.adj[b], a)
}

func (g *Graph) RemoveEdge(a, b Node) {
	for i, node := range g.adj[a] {
		if node == b {
			g.adj[a] = append(g.adj[a][:i], g.adj[a][i+1:]...)
			break
		}
	}
	for i, node := range g.adj[b] {
		if node == a {
			g.adj[b] = append(g.adj[b][:i], g.adj[b][i+1:]...)
			break
		}
	}
}

func (g *Graph) Copy() *Graph {
	ng := NewGraph()
	for _, node := range g.nodes {
		ng.nodes = append(ng.nodes, node)
		ng.adj[node] = append(ng.adj[node], g.adj[node]...)
	}
	return ng
}

func (g *Graph) Reachable(start Node) []Node {
	visited := make(map[Node]bool)
	q := []Node{start}
	visited[start] = true
	for len(q) > 0 {
		node := q[0]
		q = q[1:]
		for _, conn := range g.adj[node] {
			if visited[conn] {
				continue
			}
			visited[conn] = true
			q = append(q, conn)
		}
	}

	reachable := make([]Node, 0, len(visited))
	for node := range visited {
		reachable = append(reachable, node)
	}
	return reachable
}

type ShortestPathEntry struct {
	node Node
	path []Node
}

func (g *Graph) ShortestPath(a, b Node) []Node {
	visited := make(map[Node]bool)
	q := []ShortestPathEntry{{a, []Node{}}}
	visited[a] = true
	for len(q) > 0 {
		e := q[0]
		q = q[1:]
		if e.node == b {
			return e.path
		}
		for _, conn := range g.adj[e.node] {
			if visited[conn] {
				continue
			}

			visited[conn] = true
			q = append(q, ShortestPathEntry{conn, append(e.path, conn)})
		}
	}
	return nil
}

type Pair struct {
	a Node
	b Node
}

type KV struct {
	Key   Pair
	Value int
}

func TopK(m map[Pair]int, k int) []Pair {
	kvs := make([]KV, 0, len(m))
	for key, value := range m {
		kvs = append(kvs, KV{key, value})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})

	top := make([]Pair, 0, k)
	for i := 0; i < k; i++ {
		top = append(top, kvs[i].Key)
	}
	return top
}

func MinimumCut(g *Graph, iter int, k int) int {
	done := false
	var reachable []Node

	for !done {
		crossings := make(map[Pair]int)
		for i := 0; i < iter; i++ {
			from := g.RandomNode()
			to := g.RandomNode()
			if from == to {
				continue
			}

			path := g.ShortestPath(from, to)
			for i := 0; i < len(path)-1; i++ {
				crossings[Pair{path[i], path[i+1]}]++
			}
		}

		topK := TopK(crossings, k)
		ng := g.Copy()
		for _, pair := range topK {
			ng.RemoveEdge(pair.a, pair.b)
		}
		reachable = ng.Reachable(g.nodes[0])

		if len(reachable) != len(g.nodes) {
			done = true
		}
	}
	return len(reachable)
}

func part1() {
	lines := ReadLines("input.txt")
	g := NewGraph()
	for _, line := range lines {
		// line: "a: b c d"
		parts := strings.Split(line, ":")
		a := parts[0]
		right := strings.Split(strings.TrimSpace(parts[1]), " ")
		for _, b := range right {
			g.AddEdge(Node(a), Node(b))
		}
	}

	r := MinimumCut(g, 20, 3)
	nr := len(g.nodes) - r
	fmt.Println(r * nr)
}

func main() {
	start := time.Now()
	part1()
	fmt.Println("Part 1 finished in:", time.Since(start))
}
