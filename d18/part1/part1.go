package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ts4z/aoc2023/aoc"
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
	Color string
}

var lineRE = regexp.MustCompile(`([UDLR]) (\d+) \(#([0-9a-f]{6})\)`)

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
			string(matches[3]),
		}
		ils = append(ils, il)
	}
	return ils
}

func FindDimensions(ils []InputLine) (bottomRight matrix.Address, start matrix.Address) {
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
	pos = matrix.Address{Row: -lowestRow, Column: -lowestColumn}

	for i := len(ils) - 1; i >= 0; i-- {
		pos.Go(ils[i].RelativeAddress.Negate())
	}

	start = pos
	return
}

func findSomeInsidePoint(m *matrix.Matrix[byte]) matrix.Address {
	nRows := m.Rows()
	nCols := m.Columns()
	firstRow := nRows / 2

	for rowOffset := 0; rowOffset < nRows; rowOffset++ {
		rowN := (rowOffset + firstRow) % nRows

		for j := 0; j < nCols-1; j++ {
			maybeEdge := matrix.Address{Row: rowN, Column: j}
			maybeInside := matrix.Address{Row: rowN, Column: j + 1}

			if m.GetAddress(maybeEdge) == '#' && m.GetAddress(maybeInside) == '.' {
				return maybeInside
			}
		}
	}

	log.Fatalf("findSomeInsidePoint heuristic failed")
	return matrix.Address{}
}

func adjacentAddresses(addr matrix.Address) []matrix.Address {
	return []matrix.Address{
		{Row: addr.Row - 1, Column: addr.Column + 0},
		{Row: addr.Row + 1, Column: addr.Column + 0},
		{Row: addr.Row + 0, Column: addr.Column - 1},
		{Row: addr.Row + 0, Column: addr.Column + 1},
	}
}

func floodFillInside(m *matrix.Matrix[byte]) {
	stack := []matrix.Address{findSomeInsidePoint(m)}

	for len(stack) > 0 {
		at := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if m.GetAddress(at) == '.' {
			*m.AtAddress(at) = '#'
			stack = append(stack, adjacentAddresses(at)...)
		}
	}
}

func applyInputLinesToMatrix(m *matrix.Matrix[byte], ils []InputLine, start matrix.Address) {
	at := start
	*m.AtAddress(at) = '#'
	for _, il := range ils {
		ra := il.RelativeAddress
		log.Printf("relative address %v", ra)
		for ra.Amount > 0 {
			one := ra
			one.Amount = 1
			at = at.Go(one)
			log.Printf("set at %v", at)
			*m.AtAddress(at) = '#'
			ra.Amount--
		}
	}
}

func main() {
	ils := parseInput()
	bottomRight, start := FindDimensions(ils)
	log.Printf("bottomRight=%s start=%s", bottomRight.String(), start.String())

	m := matrix.New[byte](bottomRight.Row+1, bottomRight.Column+1)
	m.ForEachAddress(func(addr matrix.Address) {
		*m.AtAddress(addr) = '.'
	})
	applyInputLinesToMatrix(m, ils, start)

	aoc.PrintByteMatrix(m.GetRawMatrix())

	floodFillInside(m)

	vol := 0
	m.ForEachAddress(func(addr matrix.Address) {
		if m.Get(addr.Row, addr.Column) == '#' {
			vol++
		}
	})

	aoc.PrintByteMatrix(m.GetRawMatrix())

	fmt.Printf("volume %d\n", vol)
}
