package cc

// package cc has functions for analyzing "camel cards".

import (
	"log"
)

type HandType int

const (
	Nada = iota
	OnePair
	TwoPair
	Trips
	FullHouse
	Quads
	Quints
)

var cardToOrder = map[uint8]int{}

func init() {
	cardToOrder['K'] = 13
	cardToOrder['Q'] = 12
	cardToOrder['J'] = 11
	cardToOrder['T'] = 10
	cardToOrder['A'] = 14 // is this 1 or 14?
	for i := 2; i < 10; i++ {
		cardToOrder[uint8(i+'0')] = i
	}
}

func GetHandType(hand string) HandType {
	if len(hand) < 5 {
		log.Fatalf("hands have 5 cards, not %d (%q)", len(hand), hand)
	}
	if len(hand) != 5 {
		hand = hand[:5]
	}
	byCard := [15]int{}
	for i := 0; i < len(hand); i++ {
		byCard[cardToOrder[hand[i]]]++
	}
	var byMatches [6]int
	for _, v := range byCard {
		byMatches[v]++
	}
	if byMatches[5] != 0 {
		return Quints
	}
	if byMatches[4] != 0 {
		return Quads
	}
	if byMatches[3] != 0 && byMatches[2] != 0 {
		return FullHouse
	}
	if byMatches[3] != 0 {
		return Trips
	}
	if byMatches[2] > 1 {
		return TwoPair
	}
	if byMatches[2] == 1 {
		return OnePair
	}
	return Nada
}

func CompareHands(left, right string) int {
	lt := GetHandType(left)
	rt := GetHandType(right)
	if lt < rt {
		return -1
	}
	if lt > rt {
		return 1
	}
	for i := 0; i < 5; i++ {
		lch := cardToOrder[left[i]]
		rch := cardToOrder[right[i]]
		if lch != rch {
			return lch - rch
		}
	}
	log.Fatalf("can't happen, hands tied")
	return 0
}
