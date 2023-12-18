package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/ick/matrix"
)

var UDLRToDirection = map[string]matrix.Direction{
	"U": matrix.Up,
	"D": matrix.Down,
	"L": matrix.Left,
	"R": matrix.Right,
}

type InputLine struct {
	matrix.RelativeAddress
}

var lineRE = regexp.MustCompile(`([UDLR]) (\d+) \(#[0-9a-f]{6}\)`)

func parseInput() []InputLine {
	lines := ick.Must(argv.ReadChompAll())
	ils := []InputLine{}
	for _, line := range lines {
		matches := lineRE.FindSubmatch([]byte(line))
		if matches == nil {
			log.Fatalf("can't parse input line: %q", line)
		}
		il := InputLine{
			matrix.RelativeAddress{
				Direction: UDLRToDirection[string(matches[1])],
				Amount:    ick.Atoi(string(matches[2])),
			},
		}
		ils = append(ils, il)
	}
	return ils
}

func FindDimensions(ils []InputLine) (bottomRight matrix.Address, start matrix.Address) {
	// note: named return

	lowestRow := 0
	lowestColumn := 0
	highestRow := 0
	highestColumn := 0

	pos := matrix.Address{}
	for _, il := range ils {
		// was := pos
		pos = pos.Go(il.RelativeAddress)
		// log.Printf("from %s to %s, il=%+v", was, pos, il)
		if pos.Row < lowestRow {
			lowestRow = pos.Row
		}
		if pos.Column < lowestColumn {
			lowestColumn = pos.Column
		}
		if pos.Row > highestRow {
			highestRow = pos.Row
		}
		if pos.Column > highestColumn {
			highestColumn = pos.Column
		}
	}

	log.Printf("cols %d %d rows %d %d", lowestRow, highestRow, lowestColumn, highestColumn)

	finalAddress := pos
	log.Printf("finished at %v", finalAddress)
	bottomRight = matrix.Address{Row: highestRow - lowestRow, Column: highestColumn - lowestColumn}

	// start from a position that means we'll never go below 0,0 in the new 0,0
	// based matrix
	start = matrix.Address{Row: -lowestRow, Column: -lowestColumn}
	return // note: named return
}

type horizontalEdge struct {
	fromColumn int // inclusive
	toColumn   int // inclusive
}

type verticalEdge struct {
	fromRow int // inclusive
	toRow   int // inclusive
}

// https://en.wikipedia.org/wiki/Shoelace_formula
func shoelace(points []matrix.Address) int {
	area2 := 0
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		det := p1.Column*p2.Row - p1.Row*p2.Column
		log.Printf("det=%d\n", det)
		area2 += det
	}
	return area2
}

func main() {
	ils := parseInput()
	bottomRight, start := FindDimensions(ils)
	log.Printf("bottomRight=%s start=%s", bottomRight.String(), start.String())

	perimeter := 0
	points := []matrix.Address{}
	at := matrix.Address{}
	for _, il := range ils {
		perimeter += il.Amount
		np := at.Go(il.RelativeAddress)
		log.Printf("old point %v, new point %v", at, np)
		points = append(points, np)
		at = np
	}

	log.Printf("%d points", len(points))

	area := shoelace(points)
	log.Printf("perimeter %d", perimeter)
	log.Printf("2*area %d", area)

	// this answer came from AOC spoilers on reddit :-( I didn't remember the
	// shoelace formula (I can't even remember learning it), and spent a lot of
	// time coming up with complicated (and mostly wrong) ways of working this
	// out.  Guess it's time to go dig out the linear algebra book.
	fmt.Printf("%d\n", area/2+perimeter/2+1)
}
