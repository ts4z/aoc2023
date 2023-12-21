package main

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/aoc"
)

type Direction struct {
	Name           string
	Step           func(Position) Position
	NextDirections []*Direction
}

// var North = &Direction{"North", GoNorth, nil}
// var South = &Direction{"South", GoSouth, nil}
// var East = &Direction{"East", GoEast, nil}
// var West = &Direction{"West", GoWest, nil}

type Position struct {
	Row       int
	Column    int
	StepsLeft int
}

func (p Position) String() string {
	return fmt.Sprintf("Position{%d,%d, StepsLeft:%d}", p.Row, p.Column, p.StepsLeft)
}

func GoNorth(p *Position) *Position {
	return &Position{p.Row - 1, p.Column, p.StepsLeft - 1}
}
func GoWest(p *Position) *Position {
	return &Position{p.Row, p.Column - 1, p.StepsLeft - 1}
}
func GoEast(p *Position) *Position {
	return &Position{p.Row, p.Column + 1, p.StepsLeft - 1}
}
func GoSouth(p *Position) *Position {
	return &Position{p.Row + 1, p.Column, p.StepsLeft - 1}
}

func findStart(g *Grid[byte]) *Position {
	for i := range g.grid {
		for j := range g.grid[i] {
			at := &Position{Row: i, Column: j}
			if g.At(at) == 'S' {
				return at
			}
		}
	}
	log.Fatalf("can't find start position")
	return nil
}

type Grid[T any] struct {
	grid   [][]T
	maxRow int
	maxCol int
}

func NewGrid[T any](a [][]T) *Grid[T] {
	return &Grid[T]{
		grid:   a,
		maxRow: len(a),
		maxCol: len(a[0]),
	}
}

func NewGridSizedLike[T any, U any](g *Grid[U]) *Grid[T] {
	return &Grid[T]{
		grid:   aoc.NewMatrix[T](g.maxRow, g.maxCol),
		maxRow: g.maxRow,
		maxCol: g.maxCol,
	}
}

func (g *Grid[T]) At(p *Position) T {
	return g.grid[p.Row][p.Column]
}

func (g *Grid[T]) Ref(p *Position) *T {
	return &g.grid[p.Row][p.Column]
}

func (g *Grid[T]) Dimensions() (int, int) {
	return len(g.grid), len(g.grid[0])
}

func (g *Grid[T]) InBounds(p *Position) bool {
	return p.Row >= 0 && p.Row < g.maxRow && p.Column >= 0 && p.Column < g.maxCol
}

func main() {
	grid := NewGrid(aoc.ReadInputAsByteMatrix())
	start := findStart(grid)
	start.StepsLeft = 64
	visited := NewGridSizedLike[bool](grid)
	log.Printf("start at %v", start)
	reachable := 0
	steppers := []func(*Position) *Position{GoNorth, GoEast, GoSouth, GoWest}

	todo := make(chan *Position, 100000)
	todo <- start

	for len(todo) > 0 {
		here := <-todo
		log.Printf("at %v", here)

		// Reject anyplace we found another path to; our search is so this is
		// quite possible
		if visited.At(here) {
			continue
		}

		// If we start on a black square, make sure we only count as reachable if
		// we're on a black square
		if here.StepsLeft&1 == 0 {
			reachable++
		}
		// Ignore other ways to get here; because this is BFS, we have already
		// found the shortest path here
		*visited.Ref(here) = true

		// If we are out of steps, we do not generate more steps
		if here.StepsLeft == 0 {
			continue
		}

		for _, step := range steppers {

			next := step(here)

			if !grid.InBounds(next) || grid.At(next) == '#' {
				// invalid step
				continue
			}

			if visited.At(next) {
				// we've already been here, no need to do it again
				continue
			}

			todo <- next
		}
	}

	aoc.PrintBoolMatrix(visited.grid)
	fmt.Printf("reachable %d\n", reachable)
}
