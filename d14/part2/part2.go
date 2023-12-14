package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ts4z/aoc2023/argv"
)

func PrintByteMatrix(a [][]byte) {
	PrintByteMatrixTo(os.Stdout, a)
}

func PrintByteMatrixTo(w io.Writer, a [][]byte) {
	for _, line := range a {
		for _, ch := range line {
			fmt.Fprintf(w, "%c", ch)
		}
		fmt.Fprintf(w, "\n")
	}
}

type position struct{ i, j int }

type stepper func(position) position

func rollRock(a [][]byte, from position, step stepper) {
	to := step(from)
	if to.i < 0 || to.j < 0 || to.i >= len(a) || to.j >= len(a[0]) {
		a[from.i][from.j] = 'O'
		return
	}
	if a[to.i][to.j] != '.' {
		a[from.i][from.j] = 'O'
		return
	}

	a[from.i][from.j] = '.'
	rollRock(a, to, step)
}

func stepNorth(at position) position {
	return position{at.i - 1, at.j}
}

func stepSouth(at position) position {
	return position{at.i + 1, at.j}
}

func stepEast(at position) position {
	return position{at.i, at.j + 1}
}

func stepWest(at position) position {
	return position{at.i, at.j - 1}
}

func tiltTowardsNorth(a [][]byte) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] == 'O' {
				rollRock(a, position{i, j}, stepNorth)
			}
		}
	}
}

func tiltTowardsSouth(a [][]byte) {
	for i := len(a) - 1; i >= 0; i-- {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] == 'O' {
				rollRock(a, position{i, j}, stepSouth)
			}
		}
	}
}

func tiltTowardsEast(a [][]byte) {
	for j := len(a[0]) - 1; j >= 0; j-- {
		for i := 0; i < len(a); i++ {
			if a[i][j] == 'O' {
				rollRock(a, position{i, j}, stepEast)
			}
		}
	}
}

func tiltTowardsWest(a [][]byte) {
	for j := 0; j < len(a[0]); j++ {
		for i := 0; i < len(a); i++ {
			if a[i][j] == 'O' {
				rollRock(a, position{i, j}, stepWest)
			}
		}
	}
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

	const maxCycles = 1000000000
	cycleCount := 0
	cache := map[string]int{}
	cycleLength := -1

	spinCycle := func() {
		tiltTowardsNorth(a)
		tiltTowardsWest(a)
		tiltTowardsSouth(a)
		tiltTowardsEast(a)
	}

	for {
		cycleCount++

		spinCycle()

		if cycleCount&0xFFFF == 0xFFFF {
			fmt.Printf("at cycle %d\n", cycleCount)
			PrintByteMatrix(a)
			fmt.Printf("\n")
		}

		sb := &strings.Builder{}
		PrintByteMatrixTo(sb, a)
		asString := sb.String()
		if previousCycles, ok := cache[asString]; ok {
			fmt.Printf("cycle repeats: cycles=%d looks like %d\n", cycleCount, previousCycles)
			cycleLength = cycleCount - previousCycles
			fmt.Printf("cycle length is %d\n", cycleLength)
			break
		}
		cache[asString] = cycleCount
	}

	remaining := (maxCycles - cycleCount) % cycleLength

	// now, we have probably seen these cycles before too, but we'll
	// just step through until we've done the right number, it's easier
	for remaining > 0 {
		remaining--

		spinCycle()
	}

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
