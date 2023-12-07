package ccw

// package ccw has functions for analyzing "camel cards"
// where 'J' is a joker and wild.

import (
	"log"
)

type HandType int

const (
	// order is important
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
	cardToOrder['J'] = 0 // jokers are weakest card

	cardToOrder['A'] = 14
	cardToOrder['K'] = 13
	cardToOrder['Q'] = 12
	// skip 11 because I have a soul
	cardToOrder['T'] = 10
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
	jokers := 0
	for i := 0; i < len(hand); i++ {
		if hand[i] != 'J' {
			byCard[cardToOrder[hand[i]]]++
		} else {
			jokers++
		}
	}
	var byMatches [6]int
	for _, v := range byCard {
		byMatches[v]++
	}
	if jokers > 0 {
		for i := 4; i > 0; i-- {
			if byMatches[i] > 0 {
				byMatches[i]--
				byMatches[i+jokers] += jokers
				break
			}
		}
	}
	// five jokers don't match anything above, so we must
	// special-case them now; they make five of a kind
	if byMatches[5] != 0 || jokers == 5 {
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
	if lt != rt {
		return int(lt - rt)
	}
	for i := 0; i < 5; i++ {
		lch := cardToOrder[left[i]]
		rch := cardToOrder[right[i]]
		if lch == rch {
			continue
		}
		if lch != rch {
			return lch - rch
		}
	}
	log.Fatalf("can't happen, hands tied")
	return 0
}
