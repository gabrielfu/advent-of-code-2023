package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var CardStrengths = map[rune]int64{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

type Type int

const (
	Unknown Type = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func (t Type) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case FiveOfAKind:
		return "Five of a Kind"
	}
	return "Unknown"
}

var JOKER_MODE = false

type Hand struct {
	Cards        string
	Type         Type
	Bid          int
	CardsJokered string
}

func DetermineType(cards string) (Type, error) {
	counts := make(map[rune]int)
	for _, c := range cards {
		counts[c]++
	}

	var max int
	for _, v := range counts {
		if v > max {
			max = v
		}
	}

	switch max {
	case 1:
		return HighCard, nil
	case 2:
		if len(counts) == 4 {
			return OnePair, nil
		}
		return TwoPair, nil
	case 3:
		if len(counts) == 3 {
			return ThreeOfAKind, nil
		}
		return FullHouse, nil
	case 4:
		return FourOfAKind, nil
	case 5:
		return FiveOfAKind, nil
	}
	return Unknown, fmt.Errorf("unknown hand type: %s", cards)
}

func NewHand(c string, bid int) (*Hand, error) {
	t, err := DetermineType(c)
	if err != nil {
		return nil, err
	}
	return &Hand{c, t, bid, c}, nil
}

func (h *Hand) String() string {
	return fmt.Sprintf("%s (jokered=%s, type=%s, bid=%d)", h.Cards, h.CardsJokered, h.Type, h.Bid)
}

func (h *Hand) Strength() int64 {
	strength := int64(math.Pow(16, float64(5))) * int64(h.Type)
	for i, c := range h.Cards {
		s := CardStrengths[c]
		if JOKER_MODE && c == 'J' {
			s = 1
		}
		strength += int64(math.Pow(16, float64(4-i))) * s
	}
	return strength
}

func nextCartesian(a []string, r int) func() []string {
	p := make([]string, r)
	x := make([]int, len(p))
	return func() []string {
		p := p[:len(x)]
		for i, xi := range x {
			p[i] = a[xi]
		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(a) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return p
	}
}

var JokerPattern = regexp.MustCompile("J")

func (h *Hand) UseJokers() {
	indices := JokerPattern.FindAllStringIndex(h.Cards, -1)

	nc := nextCartesian(
		[]string{"A", "K", "Q", "T", "9", "8", "7", "6", "5", "4", "3", "2"},
		len(indices),
	)
	for {
		sub := nc()
		if len(sub) == 0 {
			break
		}

		cards := h.Cards
		for i, index := range indices {
			cards = cards[:index[0]] + sub[i] + cards[index[1]:]
		}
		t, err := DetermineType(cards)
		if err != nil {
			panic(err)
		}
		new := &Hand{cards, t, h.Bid, cards}
		if new.Strength() > h.Strength() {
			h.CardsJokered = cards
			h.Type = t
		}
	}
}

func ReadLines(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

func part1() {
	JOKER_MODE = false
	lines := ReadLines("input.txt")
	var hands []*Hand
	for _, line := range lines {
		split := strings.Split(line, " ")

		bid, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}

		hand, err := NewHand(split[0], bid)
		if err != nil {
			panic(err)
		}
		hands = append(hands, hand)
	}

	// sort by ascending order
	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Strength() < hands[j].Strength()
	})
	winnings := 0
	for i, hand := range hands {
		winnings += hand.Bid * (i + 1)
	}
	println(winnings)
}

func part2() {
	JOKER_MODE = true
	lines := ReadLines("input.txt")
	var hands []*Hand
	for _, line := range lines {
		split := strings.Split(line, " ")

		bid, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}

		hand, err := NewHand(split[0], bid)
		if err != nil {
			panic(err)
		}
		hands = append(hands, hand)
	}

	for _, hand := range hands {
		hand.UseJokers()
	}

	// sort by ascending order
	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Strength() < hands[j].Strength()
	})
	winnings := 0
	for i, hand := range hands {
		winnings += hand.Bid * (i + 1)
	}
	println(winnings)
}

func main() {
	part1()
	part2()
}
