package main

import (
	"log"
	"sort"
	"strconv"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/d7/cc"
	"github.com/ts4z/aoc2023/ick"
)

func main() {
	lines := ick.Must(argv.ReadChompAll())
	sort.Slice(lines, func(i, j int) bool {
		// reverse order
		return cc.CompareHands(lines[i], lines[j]) < 0
	})

	total := 0
	for i, line := range lines {
		rank := i + 1
		bid := ick.Must(strconv.Atoi(line[6:]))
		winnings := rank * bid
		total += winnings
		log.Printf("rank=%d line=%q bid=%d winnings=%d", rank, line, bid, winnings)
	}
	log.Printf("total winnings = %d", total)
}
