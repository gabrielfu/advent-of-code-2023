package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
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

type Pulse int

const (
	HighPulse Pulse = iota
	LowPulse
)

func (p Pulse) String() string {
	switch p {
	case HighPulse:
		return "high"
	case LowPulse:
		return "low"
	default:
		return "?"
	}
}

type ModuleType int

const (
	UntypedModuleType ModuleType = iota
	BroadcasterModuleType
	FlipFlopModuleType
	ConjunctionModuleType
)

type Instruction struct {
	source string
	pulse  Pulse
	dest   string
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s -%s-> %s", i.source, i.pulse, i.dest)
}

type Module interface {
	Name() string
	Destinations() []string
	Type() ModuleType
	Instructions(string, Pulse) []Instruction
}

type BroadcasterModule struct {
	name         string
	destinations []string
}

func (m *BroadcasterModule) String() string {
	return fmt.Sprintf("BroadcasterModule(%s, %v)", m.name, m.destinations)
}

func (m *BroadcasterModule) Name() string {
	return m.name
}

func (m *BroadcasterModule) Destinations() []string {
	return m.destinations
}

func (m *BroadcasterModule) Type() ModuleType {
	return BroadcasterModuleType
}

func (m *BroadcasterModule) Instructions(source string, pulse Pulse) []Instruction {
	var instructions []Instruction
	for _, destination := range m.destinations {
		instructions = append(instructions, Instruction{m.name, pulse, destination})
	}
	return instructions
}

type FlipFlopModule struct {
	name         string
	destinations []string
	turnedOn     bool
}

func (m *FlipFlopModule) String() string {
	return fmt.Sprintf("FlipFlopModule(%s, %v)", m.name, m.destinations)
}

func (m *FlipFlopModule) Name() string {
	return m.name
}

func (m *FlipFlopModule) Destinations() []string {
	return m.destinations
}

func (m *FlipFlopModule) Type() ModuleType {
	return FlipFlopModuleType
}

func (m *FlipFlopModule) Instructions(source string, pulse Pulse) []Instruction {
	var instructions []Instruction
	if pulse == LowPulse {
		before := m.IsTurnedOn()
		m.Toggle()
		if before {
			for _, destination := range m.destinations {
				instructions = append(instructions, Instruction{m.name, LowPulse, destination})
			}
		} else {
			for _, destination := range m.destinations {
				instructions = append(instructions, Instruction{m.name, HighPulse, destination})
			}
		}
	}
	return instructions
}

func (m *FlipFlopModule) IsTurnedOn() bool {
	return m.turnedOn
}

func (m *FlipFlopModule) Toggle() {
	m.turnedOn = !m.turnedOn
}

type ConjunctionModule struct {
	name         string
	destinations []string
	inputs       map[string]Pulse
}

func (m *ConjunctionModule) String() string {
	return fmt.Sprintf("ConjunctionModule(%s, %v)", m.name, m.destinations)
}

func (m *ConjunctionModule) Name() string {
	return m.name
}

func (m *ConjunctionModule) Destinations() []string {
	return m.destinations
}

func (m *ConjunctionModule) Type() ModuleType {
	return ConjunctionModuleType
}

func (m *ConjunctionModule) Instructions(source string, pulse Pulse) []Instruction {
	var instructions []Instruction
	m.inputs[source] = pulse
	if m.AllInputsAreHigh() {
		for _, destination := range m.destinations {
			instructions = append(instructions, Instruction{m.name, LowPulse, destination})
		}
	} else {
		for _, destination := range m.destinations {
			instructions = append(instructions, Instruction{m.name, HighPulse, destination})
		}
	}
	return instructions
}

func (m *ConjunctionModule) AllInputsAreHigh() bool {
	for _, pulse := range m.inputs {
		if pulse == LowPulse {
			return false
		}
	}
	return true
}

func (m *ConjunctionModule) AllInputsAreLow() bool {
	for _, pulse := range m.inputs {
		if pulse == HighPulse {
			return false
		}
	}
	return true
}

func (m *ConjunctionModule) AddInput(name string) {
	m.inputs[name] = LowPulse
}

var pattern = regexp.MustCompile(`^([%&]?)([a-z]+) -> ([a-z, ]+)$`)

func ParseLine(line string) Module {
	matches := pattern.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	var module Module
	name := matches[2]
	dests := strings.Split(matches[3], ", ")
	switch matches[1] {
	case "%":
		module = &FlipFlopModule{name, dests, false}
	case "&":
		module = &ConjunctionModule{name, dests, make(map[string]Pulse)}
	case "":
		if name == "broadcaster" {
			module = &BroadcasterModule{name, dests}
		} else {
			panic("unknown module type")
		}
	default:
		panic("unknown module type")
	}

	return module
}

type Registry struct {
	r       map[string]Module
	rxCount int
}

func (r *Registry) Items() map[string]Module {
	return r.r
}

func NewRegistry(lines []string) *Registry {
	var registry = make(map[string]Module)
	for _, line := range lines {
		module := ParseLine(line)
		registry[module.Name()] = module
	}

	for _, module := range registry {
		for _, destination := range module.Destinations() {
			destModule, ok := registry[destination]
			if !ok {
				continue
			}
			if destModule.Type() == ConjunctionModuleType {
				destModule.(*ConjunctionModule).AddInput(module.Name())
			}
		}
	}
	return &Registry{registry, 0}
}

func (r *Registry) PressButton(hook *Hook, i int) (int, int) {
	instructions := []Instruction{{"button", LowPulse, "broadcaster"}}
	var low, high int
	for len(instructions) > 0 {
		instruction := instructions[0]
		instructions = instructions[1:]

		if hook != nil && instruction.pulse == LowPulse && slices.Contains(hook.hooks, instruction.dest) {
			if hook.loops[instruction.dest] == 0 {
				hook.loops[instruction.dest] = i
				index := slices.Index(hook.hooks, instruction.dest)
				hook.hooks = append(hook.hooks[:index], hook.hooks[index+1:]...)
			}
		}

		if instruction.pulse == LowPulse {
			low++
		} else {
			high++
		}
		m, ok := r.r[instruction.dest]
		if !ok {
			continue
		}
		instructions = append(instructions, m.Instructions(instruction.source, instruction.pulse)...)
	}
	return low, high
}

func (r *Registry) IsInitialState() bool {
	for _, module := range r.Items() {
		switch module.Type() {
		case FlipFlopModuleType:
			if module.(*FlipFlopModule).IsTurnedOn() {
				return false
			}
		case ConjunctionModuleType:
			if !module.(*ConjunctionModule).AllInputsAreLow() {
				return false
			}
		default:
		}
	}
	return true
}

func part1() {
	lines := ReadLines("input.txt")
	registry := NewRegistry(lines)

	N := 1000
	var period int
	var low, high int
	for {
		l, h := registry.PressButton(nil, 0)
		low += l
		high += h
		period++
		if registry.IsInitialState() || period >= N {
			break
		}
	}

	div, mod := N/period, N%period
	low *= div
	high *= div

	for i := 0; i < mod; i++ {
		l, h := registry.PressButton(nil, 0)
		low += l
		high += h
	}

	fmt.Println(low * high)
}

type Hook struct {
	hooks []string
	loops map[string]int
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
	registry := NewRegistry(lines)

	// only &vr connects to rx
	// the inputs of vr are all conjunction modules bm, cl, tn & dr
	// thus, we find the cycle of each conjunction module sending a low pulse
	// to vr. When these 4 modules send a low pulse to vr together,
	// vr will send a low pulse to rx. So, we just need to find the lcm of these cycles

	hook := &Hook{
		hooks: []string{"bm", "cl", "tn", "dr"},
		loops: map[string]int{
			"bm": 0,
			"cl": 0,
			"tn": 0,
			"dr": 0,
		},
	}
	var i int
	for {
		i++
		registry.PressButton(hook, i)
		if len(hook.hooks) == 0 {
			break
		}
	}

	// lcm of hook.loops values
	var values []int
	for _, value := range hook.loops {
		values = append(values, value)
	}
	lcm := LCM(values[0], values[1], values...)
	fmt.Println(lcm)
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
