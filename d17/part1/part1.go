package main

// this has an off-by-one error I haven't found, and on both the example input
// as well as my input, the value is 1 too low.

import (
	"fmt"

	"github.com/ts4z/aoc2023/aoc"
	"github.com/ts4z/aoc2023/ick/pq"
)

type Position struct {
	Row           int
	Column        int
	LastDirection *Direction
	Straights     int
	Cost          int
}

func (pos Position) Priority() int {
	return pos.Cost
}

func (pos Position) String() string {
	return fmt.Sprintf("Position{%d,%d,%s,%d,%d}",
		pos.Row, pos.Column, pos.LastDirection.Name, pos.Straights, pos.Cost)
}

type Direction struct {
	Name           string
	Step           func(Position) Position
	NextDirections []*Direction
}

var NoDirection = &Direction{"None", nil, []*Direction{South, East}}
var North = &Direction{"North", GoNorth, nil}
var South = &Direction{"South", GoSouth, nil}
var East = &Direction{"East", GoEast, nil}
var West = &Direction{"West", GoWest, nil}

func init() {
	North.NextDirections = []*Direction{North, East, West}
	South.NextDirections = []*Direction{South, East, West}
	East.NextDirections = []*Direction{North, South, East}
	West.NextDirections = []*Direction{North, South, West}
}

func GoNorth(p Position) Position {
	p.Row--
	return p
}
func GoWest(p Position) Position {
	p.Column--
	return p
}
func GoEast(p Position) Position {
	p.Column++
	return p
}
func GoSouth(p Position) Position {
	p.Row++
	return p
}

func CheckBounds[T any](a [][]T, pos Position) bool {
	if pos.Row < 0 {
		return false
	}
	if pos.Row >= len(a) {
		return false
	}
	if pos.Column < 0 {
		return false
	}
	if pos.Column >= len(a[0]) {
		return false
	}
	return true // ok
}

func process(a [][]int) int {

	todo := pq.New[*Position]()
	didIt := map[string]struct{}{}
	goalRow := len(a) - 1
	goalCol := len(a[0]) - 1

	// Start in the top-left without facing a direction (a special case).
	todo.Push(&Position{0, 0, NoDirection, 0, 0})

	for todo.Len() > 0 {
		pos := todo.Pop()

		// log.Printf("I'm at %v", pos)

		key := pos.String()
		if _, ok := didIt[key]; ok {
			continue
		}
		didIt[key] = struct{}{}

		if pos.Row == goalRow && pos.Column == goalCol {
			return pos.Cost
		}

		for _, dir := range pos.LastDirection.NextDirections {
			nextPos := dir.Step(*pos)
			if !CheckBounds(a, nextPos) {
				// log.Printf("out of bounds: %v", nextPos)
				continue
			}
			if pos.LastDirection != dir {
				nextPos.Straights = 1
			} else {
				nextPos.Straights = pos.Straights + 1
				if nextPos.Straights > 3 {
					// log.Printf("too many straights: %v", nextPos)
					continue
				}
			}

			nextPos.LastDirection = dir
			nextPos.Cost = pos.Cost + a[nextPos.Row][nextPos.Column]

			todo.Push(&nextPos)
		}
	}

	fmt.Printf("no answer found!\n")
	return -1
}

type PriorityInt int

func (pi PriorityInt) Priority() int {
	return int(pi)
}

func main() {
	ab := aoc.ReadInputAsByteMatrix()

	a := make([][]int, len(ab))
	for i, row := range ab {
		a[i] = make([]int, len(ab[i]))
		for j, b := range row {
			a[i][j] = int(b - '0')
		}
	}

	// initialPosition := Position{0, 0, NoDirection, 0, 0}
	v := process(a)

	fmt.Printf("%d\n", v)
}
