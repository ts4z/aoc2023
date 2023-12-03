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
			// log.Printf("update %d,%d? yes", i, j)
			c := &adjacency[i][j]
			c.count++
			c.product *= number

			// log.Printf("updated %d,%d: %+v\n", i, j, adjacency[i][j])
		} else {
			// log.Printf("update %d,%d? no, oob\n", i, j)
		}
	}

	// Walk the array looking for numbers.
	// Numbers are horizontal.
	// If we find one, add a count and a number to all adjacent cells.
	//
	// We do this left-to-right so that we handle each number once,
	// no matter how many cells it's in.
	for i, line := range lines {
		working := 0
		inNumber := false
		firstDigitAt := -1

		finishNumber := func(nextNonDigitAt int) {
			// log.Printf("finish number at %d,%d", i, nextNonDigitAt)
			number := working

			// No need to update inside the number itself, but on the
			// line with the number, do left and right
			maybeSet(i, firstDigitAt-1, number)
			maybeSet(i, nextNonDigitAt, number)
			for j := firstDigitAt - 1; j < nextNonDigitAt+1; j++ {
				// update all cells above the number
				maybeSet(i-1, j, number)
				// and below the number
				maybeSet(i+1, j, number)
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
					// log.Printf("j=%d now in, firstDigitAt %d", j, firstDigitAt)
				}
				working *= 10
				working += int(ch - '0')
			} else if inNumber {
				// we were in a number, now we're not
				// log.Printf("j=%d now out", j)
				finishNumber(j)
			}
		}

		if inNumber {
			finishNumber(len(line))
		}
	}

	// Walk the array AGAIN, and look for gears.
	// If a gear is found and it is adjacent to two numbers,
	// collect the product of those numbers in the sum
	sum := 0
	for i, line := range lines {
		for j, ch := range line {
			c := &adjacency[i][j]
			if ch == '*' {
				// log.Printf("gear %d,%d\n", i, j)
				if c.count == 2 {
					// log.Printf("count 2 at gear %d,%d\n", i, j)
					sum += c.product
				} else {
					// log.Printf("count %d, not 2\n", c.count)
				}
			}
		}
	}

	fmt.Printf("sum %d\n", sum)
}
