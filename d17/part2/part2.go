package main

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/aoc"
	"github.com/ts4z/aoc2023/ick/pq"
)

type Position struct {
	Row           int
	Column        int
	LastDirection *Direction
	Straights     int
	CostToGetHere int
	Previous      *Position
}

func (pos Position) Priority() int {
	return pos.CostToGetHere
}

func (pos Position) String() string {
	return fmt.Sprintf("Position{%d,%d,%s,%d,%d}",
		pos.Row, pos.Column, pos.LastDirection.Name, pos.Straights, pos.CostToGetHere)
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

	// Start in the top-left without facing a direction (a special case).  But we
	// claim 5 straights so we are allowed to change direction immediately.
	todo.Push(&Position{0, 0, NoDirection, 5, 0, nil})

	for todo.Len() > 0 {
		pos := todo.Pop()

		key := pos.String()
		if _, ok := didIt[key]; ok {
			continue
		}
		didIt[key] = struct{}{}

		insertNext := func(nextPos *Position) {
			nextPos.CostToGetHere = pos.CostToGetHere + a[nextPos.Row][nextPos.Column]
			nextPos.Previous = pos

			todo.Push(nextPos)
		}

		if pos.Straights < 4 {
			// We need to go at least 4 straights from this state; we can't even get
			// to the goal until we do.

			nextPos := pos.LastDirection.Step(*pos)

			if !CheckBounds(a, nextPos) {
				// log.Printf("out of bounds: %v", nextPos)
				continue
			}

			nextPos.LastDirection = pos.LastDirection
			nextPos.Straights = pos.Straights + 1
			insertNext(&nextPos)

		} else {

			if pos.Row == goalRow && pos.Column == goalCol {
				log.Printf("answer!")
				for cur := pos; cur != nil; cur = cur.Previous {
					log.Printf("  %v", cur.String())
				}
				return pos.CostToGetHere
			}

			for _, dir := range pos.LastDirection.NextDirections {
				nextPos := dir.Step(*pos)
				if !CheckBounds(a, nextPos) {
					// log.Printf("out of bounds: %v", nextPos)
					continue
				}

				if pos.LastDirection != dir {
					if pos.Straights < 4 {
						// log.Printf("too FEW straights: %v", nextPos)
						continue
					}

					// OK, we can change direction now.
					nextPos.Straights = 1
				} else {
					nextPos.Straights = pos.Straights + 1
					if nextPos.Straights > 10 {
						// log.Printf("too many straights: %v", nextPos)
						continue
					}
				}

				nextPos.LastDirection = dir
				insertNext(&nextPos)
			}
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

	v := process(a)

	fmt.Printf("%d\n", v)
}
