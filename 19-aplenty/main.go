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

type Cmp string

const (
	GT  Cmp = ">"
	LT  Cmp = "<"
	GTE Cmp = ">="
	LTE Cmp = "<="
)

func (c Cmp) Opposite() Cmp {
	switch c {
	case GT:
		return LTE
	case LT:
		return GTE
	case GTE:
		return LT
	case LTE:
		return GT
	default:
		panic("invalid Cmp")
	}
}

type Condition struct {
	empty bool
	cat   string
	cmp   Cmp
	val   int
}

func (c Condition) String() string {
	if c.empty {
		return "Condition(empty)"
	}
	return fmt.Sprintf("Condition(%s%s%d)", c.cat, c.cmp, c.val)
}

func (c Condition) Copy() Condition {
	return Condition{
		empty: c.empty,
		cat:   c.cat,
		cmp:   c.cmp,
		val:   c.val,
	}
}

func (c Condition) Opposite() Condition {
	if c.empty {
		return c
	}
	return Condition{
		empty: false,
		cat:   c.cat,
		cmp:   c.cmp.Opposite(),
		val:   c.val,
	}
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
	case GT:
		return value > c.val
	case LT:
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

func (w Workflow) Len() int {
	return len(w.rules)
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
			cmp:   Cmp(matches[2]),
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

func ParseLines(lines []string) (map[string]Workflow, []Part) {
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
	return workflows, parts
}

func part1() {
	lines := ReadLines("input.txt")
	workflows, parts := ParseLines(lines)

	total := 0
	for _, p := range parts {
		result := RunWorkflows(workflows, p, "in")
		if result == "A" {
			total += p.TotalRating()
		}
	}
	println(total)
}

type Conditions []Condition

func (c Conditions) Copy() Conditions {
	var copy Conditions
	for _, cond := range c {
		copy = append(copy, cond.Copy())
	}
	return copy
}

type Range struct {
	min, max int // inclusive
}

func NewRange() *Range {
	return &Range{1, 4000}
}

func (r Range) String() string {
	return fmt.Sprintf("Range(%d, %d)", r.min, r.max)
}

func (c Conditions) Consolidate() map[string]*Range {
	ranges := make(map[string]*Range)
	ranges["x"] = NewRange()
	ranges["m"] = NewRange()
	ranges["a"] = NewRange()
	ranges["s"] = NewRange()
	for _, cond := range c {
		if cond.empty {
			continue
		}
		r := ranges[cond.cat]
		switch cond.cmp {
		case GT:
			r.min = max(r.min, cond.val+1)
		case LT:
			r.max = min(r.max, cond.val-1)
		case GTE:
			r.min = max(r.min, cond.val)
		case LTE:
			r.max = min(r.max, cond.val)
		default:
			panic("invalid Cmp")
		}
	}
	return ranges
}

type Accepted struct {
	data []Conditions
}

func (a *Accepted) Add(c Conditions) {
	a.data = append(a.data, c)
}

func traverse(workflows map[string]Workflow, wname string, r int, conditions Conditions, accepted *Accepted) {
	w := workflows[wname]
	rule := w.rules[r]
	if rule.cond.empty {
		if rule.dest == "A" {
			// fmt.Println("Accepted!", conditions)
			accepted.Add(conditions)
		} else if rule.dest == "R" {
			return
		} else {
			traverse(workflows, rule.dest, 0, conditions, accepted)
		}
	} else {
		if rule.dest == "A" {
			// fmt.Println("Accepted!", append(conditions.Copy(), rule.cond))
			accepted.Add(append(conditions.Copy(), rule.cond))
		} else if rule.dest != "R" {
			traverse(workflows, rule.dest, 0, append(conditions.Copy(), rule.cond), accepted)
		}

		if r+1 < w.Len() {
			traverse(workflows, wname, r+1, append(conditions.Copy(), rule.cond.Opposite()), accepted)
		}
	}
}

func Solve2(workflows map[string]Workflow, start string, conditions Conditions) int {
	accepted := &Accepted{}
	traverse(workflows, start, 0, conditions, accepted)
	total := 0
	for _, c := range accepted.data {
		ranges := c.Consolidate()
		if len(ranges) == 0 {
			continue
		}
		comb := 1
		for _, r := range ranges {
			comb *= (r.max - r.min + 1)
		}
		total += comb
	}
	return total
}

func part2() {
	lines := ReadLines("input.txt")
	workflows, _ := ParseLines(lines)
	total := Solve2(workflows, "in", Conditions{})
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
