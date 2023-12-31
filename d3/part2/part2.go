package main

import (
	"fmt"
	"log" // all kids love log

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

type cell struct {
	count   int
	product int
}

func main() {
	lines := ick.Must(argv.ReadChompAll())
	lineLen := len(lines[0])

	// Double-check the input.
	for i, line := range lines {
		if lineLen != len(line) {
			log.Fatalf("line length was %v, not %v, at line %v", len(line), lineLen, i)
		}
	}

	// "adjacency" misnamed from part 1.  Ultimately, this represents
	// the product of this cell if this cell happend to contain a *.
	//
	// We avoid doing update here if the cell did not contain a * below,
	// so this isn't valid for the whole grid, just the interesting cells.
	adjacency := ick.New2DArrayWithDefault[cell](len(lines), lineLen,
		cell{count: 0, product: 1})

	// Adjust a cell if the access is in bound and might be interesting.
	maybeSet := func(i, j int, number int) {
		if i >= 0 && i < len(adjacency) &&
			j >= 0 && j < len(adjacency[i]) &&
			lines[i][j] == '*' {
			c := &adjacency[i][j]
			c.count++
			c.product *= number
		}
	}

	// Walk the array looking for numbers.  Numbers are horizontal.  If we find
	// one, add a count and a number to all adjacent cells.
	//
	// We do this left-to-right so that we handle each number once,
	// no matter how many cells it takes for the number.
	for i, line := range lines {
		working := 0
		inNumber := false
		firstDigitAt := -1

		finishNumber := func(nextNonDigitAt int) {
			// No need to update inside the number itself, but on the
			// line with the number, do left and right
			maybeSet(i, firstDigitAt-1, working)
			maybeSet(i, nextNonDigitAt, working)
			for j := firstDigitAt - 1; j < nextNonDigitAt+1; j++ {
				// update all cells above the number
				maybeSet(i-1, j, working)
				// and below the number
				maybeSet(i+1, j, working)
			}

			inNumber = false
		}

		for j, ch := range line {
			// log.Printf("j=%d", j)
			if ch >= '0' && ch <= '9' {
				// we are either starting a number or accumulating a number
				if !inNumber {
					// we are starting the number
					working = 0
					inNumber = true
					firstDigitAt = j
				}
				working *= 10
				working += int(ch - '0')
			} else if inNumber {
				// we were in a number, now we're not
				finishNumber(j)
			}
		}

		if inNumber {
			finishNumber(len(line))
		}
	}

	sum := 0
	for i, line := range lines {
		for j := range line {
			c := &adjacency[i][j]
			if c.count == 2 {
				sum += c.product
			}
		}
	}

	fmt.Printf("sum %d\n", sum)
}
