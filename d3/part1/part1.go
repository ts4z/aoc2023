package main

import (
	"fmt"
	"log" // all kids love log

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

func main() {
	lines, _ := argv.ReadChompAll()
	lineLen := len(lines[0])

	for i, line := range lines {
		if lineLen != len(line) {
			log.Fatalf("line length was %v, not %v, at line %v", len(line), lineLen, i)
		}
	}

	adjacency := ick.New2DArray[bool](len(lines), len(lines[0]))

	maybeSet := func(i, j int) {
		if i >= 0 && i < len(adjacency) &&
			j >= 0 && j < len(adjacency[i]) {
			log.Printf("set adj %d,%d", i, j)
			adjacency[i][j] = true
		}
	}

	// prepare adjacency array
	for i, line := range lines {
		for j := 0; j < lineLen; j++ {
			ch := line[j]
			log.Printf("consider %d,%d %c", i, j, ch)
			log.Printf("ch != . %v", ch != '.')
			log.Printf("ch < '0' %v", ch < '0')
			log.Printf("ch > '9' %v", ch > '9')
			if ch != '.' && (ch < '0' || ch > '9') {
				log.Printf("adjacent to %d,%d set", i, j)
				maybeSet(i-1, j-1)
				maybeSet(i-1, j)
				maybeSet(i-1, j+1)
				maybeSet(i, j-1)
				maybeSet(i, j)
				maybeSet(i, j+1)
				maybeSet(i+1, j-1)
				maybeSet(i+1, j)
				maybeSet(i+1, j+1)
			}
		}
	}

	sum := 0
	for i, line := range lines {
		working := 0
		inNumber := false
		isAdjacent := false

		finishNumber := func() {
			if isAdjacent {
				fmt.Printf("isAdjacent")
				sum += working
			}

			working = 0
			inNumber = false
			isAdjacent = false
		}

		for j, ch := range line {
			log.Printf("adjacency[%d] = %v", i, adjacency[i])
			log.Printf("adjacency[%d][%d] = %v", i, j, adjacency[i][j])
			log.Printf("j=%d", j)
			// we are either starting a number or accumulating a number
			if ch >= '0' && ch <= '9' {
				log.Printf("in number")
				inNumber = true
				working *= 10
				working += int(ch - '0')
				isAdjacent = isAdjacent || adjacency[i][j]
				log.Printf("isAdjacent is now %v", isAdjacent)
			} else if inNumber {
				finishNumber()
			} else {
				log.Printf("not in number, ch %d, %c out of range", ch, rune(ch))
			}
		}

		finishNumber()
	}

	fmt.Printf("sum %v\n", sum)
}
