package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ts4z/aoc2023/argv"
)

func PrintByteMatrix(a [][]byte) {
	w := os.Stdout
	for _, line := range a {
		for _, ch := range line {
			fmt.Fprintf(w, "%c", ch)
		}
		fmt.Fprintf(w, "\n")
	}
}

func tryToRollRockNorthTo(a [][]byte, i, j int) bool {
	if i < 0 {
		return false // caller keeps rock
	}
	if a[i][j] != '.' {
		return false // caller keeps rock
	}

	// Rock rolls to at least row,col.  See if we can roll it further north.
	if tryToRollRockNorthTo(a, i-1, j) {
		a[i][j] = '.' // caller took rock
	} else {
		a[i][j] = 'O' // we kept rock
	}
	return true // I took rock from caller
}

func main() {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	a := make([][]byte, len(lines))
	for i, line := range lines {
		a[i] = []byte(line)
	}

	fmt.Printf("input\n")
	PrintByteMatrix(a)

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] == 'O' {
				if tryToRollRockNorthTo(a, i-1, j) {
					a[i][j] = '.'
				}
			}
		}
	}

	fmt.Printf("north tilted matrix\n")
	PrintByteMatrix(a)

	tw := 0
	for i := 0; i < len(a); i++ {
		w := len(a) - i
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] == 'O' {
				tw += w
			}
		}
	}

	fmt.Printf("total weight %d\n", tw)
}
