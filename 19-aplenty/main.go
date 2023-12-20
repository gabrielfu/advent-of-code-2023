package main

import (
	"fmt"
	"os"
	"regexp"
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

type Part struct {
	x, m, a, s int
}

func (p Part) String() string {
	return fmt.Sprintf("Part(x=%d, m=%d, a=%d, s=%d)", p.x, p.m, p.a, p.s)
}

func (p Part) TotalRating() int {
	return p.x + p.m + p.a + p.s
}

var partPattern = regexp.MustCompile(`^{x=([0-9]+),m=([0-9]+),a=([0-9]+),s=([0-9]+)}$`)

func ParsePart(line string) Part {
	matches := partPattern.FindStringSubmatch(line)
	if matches == nil {
		panic("invalid part: " + line)
	}
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}
	m, err := strconv.Atoi(matches[2])
	if err != nil {
		panic(err)
	}
	a, err := strconv.Atoi(matches[3])
	if err != nil {
		panic(err)
	}
	s, err := strconv.Atoi(matches[4])
	if err != nil {
		panic(err)
	}
	return Part{x, m, a, s}
}

type Condition struct {
	empty bool
	cat   string
	cmp   string
	val   int
}

func (c Condition) String() string {
	if c.empty {
		return "Condition(empty)"
	}
	return fmt.Sprintf("Condition(%s%s%d)", c.cat, c.cmp, c.val)
}

func (c Condition) Apply(p Part) bool {
	if c.empty {
		return true
	}
	var value int
	switch c.cat {
	case "x":
		value = p.x
	case "m":
		value = p.m
	case "a":
		value = p.a
	case "s":
		value = p.s
	default:
	}
	switch c.cmp {
	case ">":
		return value > c.val
	case "<":
		return value < c.val
	default:
	}
	return true
}

type Rule struct {
	cond Condition
	dest string
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule(%v, %s)", r.cond, r.dest)
}

type Workflow struct {
	name  string
	rules []Rule
}

func (w Workflow) String() string {
	return fmt.Sprintf("Workflow(%s, rules=%v)", w.name, w.rules)
}

var rulePattern = regexp.MustCompile(`^([xmas])([<>])([0-9]+):([a-zA-Z]+)$`)

func ParseWorkflow(line string) Workflow {
	split := strings.Split(line, "{")
	name := strings.TrimSpace(split[0])
	rulePart := strings.Replace(split[1], "}", "", -1)

	rules := make([]Rule, 0)
	for _, rule := range strings.Split(rulePart, ",") {
		matches := rulePattern.FindStringSubmatch(rule)
		if matches == nil {
			rules = append(rules, Rule{Condition{empty: true}, rule})
			break
		}
		val, err := strconv.Atoi(matches[3])
		if err != nil {
			panic(err)
		}
		cond := Condition{
			cat:   matches[1],
			cmp:   matches[2],
			val:   val,
			empty: false,
		}
		rules = append(rules, Rule{cond, matches[4]})
	}
	return Workflow{name, rules}
}

// Returns the destination workflow name, or "R" / "A" for reject / accept
func (w Workflow) Apply(p Part) string {
	for _, rule := range w.rules {
		if rule.cond.Apply(p) {
			return rule.dest
		}
	}
	panic("no rule applies")
}

func RunWorkflows(ws map[string]Workflow, p Part, start string) string {
	w := ws[start]
	for {
		next := w.Apply(p)
		if next == "R" || next == "A" {
			return next
		}
		w = ws[next]
	}
}

func part1() {
	lines := ReadLines("input.txt")
	workflows := make(map[string]Workflow)

	var i int
	for _, line := range lines {
		i++
		if line == "" {
			break
		}
		w := ParseWorkflow(line)
		workflows[w.name] = w
	}

	var parts []Part
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}
		parts = append(parts, ParsePart(line))
	}

	total := 0
	for _, p := range parts {
		result := RunWorkflows(workflows, p, "in")
		if result == "A" {
			total += p.TotalRating()
		}
	}
	println(total)

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
