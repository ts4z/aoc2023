package main

/* this almost worked, then I made a mess of it.  this approach can work
   but it needs to walk the perimeter more carefully than I did, and there
   are more types of edge nodes.

   I give up, this one isn't any fun.  The clever solution recognizes more
   convenient things about the input and just works out the factors and the
   quadratic nature of the problem.
*/

import (
	"flag"
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/aoc"
)

type Direction struct {
	Name           string
	Step           func(PositionAndSteps) PositionAndSteps
	NextDirections []*Direction
}

// var North = &Direction{"North", GoNorth, nil}
// var South = &Direction{"South", GoSouth, nil}
// var East = &Direction{"East", GoEast, nil}
// var West = &Direction{"West", GoWest, nil}

type Position struct {
	Row    int
	Column int
}

type PositionAndSteps struct {
	Position
	StepsTaken int
}

func (p PositionAndSteps) String() string {
	return fmt.Sprintf("Position{%d,%d, StepsTaken:%d}", p.Row, p.Column, p.StepsTaken)
}

func (p Position) String() string {
	return fmt.Sprintf("Position{%d,%d}", p.Row, p.Column)
}

func GoNorth(p *PositionAndSteps) *PositionAndSteps {
	return &PositionAndSteps{Position{p.Row - 1, p.Column}, p.StepsTaken + 1}
}
func GoWest(p *PositionAndSteps) *PositionAndSteps {
	return &PositionAndSteps{Position{p.Row, p.Column - 1}, p.StepsTaken + 1}
}
func GoEast(p *PositionAndSteps) *PositionAndSteps {
	return &PositionAndSteps{Position{p.Row, p.Column + 1}, p.StepsTaken + 1}
}
func GoSouth(p *PositionAndSteps) *PositionAndSteps {
	return &PositionAndSteps{Position{p.Row + 1, p.Column}, p.StepsTaken + 1}
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

func (g *Grid[T]) Dimensions() (int, int) {
	return len(g.grid), len(g.grid[0])
}

func (g *Grid[T]) InBounds(p *Position) bool {
	return p.Row >= 0 && p.Row < g.maxRow && p.Column >= 0 && p.Column < g.maxCol
}

type VisitInfo struct {
	Start           *Position
	StepsToPosition map[string]int // string is Position.String()
	Reachable       int
}

func MakeVisitInfo(grid *Grid[byte], start *Position, maxSteps int) *VisitInfo {
	var visited = map[string]int{}

	// log.Printf("start at %v", start)
	// maxRow, maxCol := grid.Dimensions()
	// log.Printf("dimensions %v %v", maxRow, maxCol)
	reachable := 0
	steppers := []func(*PositionAndSteps) *PositionAndSteps{GoNorth, GoEast, GoSouth, GoWest}

	todo := make(chan *PositionAndSteps, 10000)
	todo <- &PositionAndSteps{*start, 0}

	for len(todo) > 0 {
		here := <-todo

		// Reject anyplace we found another path to; our search is so this is
		// quite possible
		if _, ok := visited[here.String()]; ok {
			log.Printf("been to %v before", here)
			continue
		}

		if (here.StepsTaken & 1) == 0 {
			reachable++
		}
		// Ignore other ways to get here; because this is BFS, we have already
		// found the shortest path here
		visited[here.String()] = here.StepsTaken
		log.Printf("been to %d places, %d in queue", len(visited), len(todo))

		// If we are out of steps, we do not generate more steps below
		if here.StepsTaken > maxSteps {
			continue
		}

		for _, step := range steppers {

			next := step(here)

			if !grid.InBounds(&next.Position) || grid.At(&next.Position) == '#' {
				// invalid step
				continue
			}

			if _, ok := visited[next.Position.String()]; ok {
				log.Printf("new step: been here before")
				// we've already been here, no need to do it again
				continue
			}

			select {
			case todo <- next:
			default:
				panic("queue full")
			}
		}
	}

	return &VisitInfo{
		Start:           start,
		StepsToPosition: visited,
		Reachable:       reachable,
	}
}

const maxSteps = 26501365

func main() {
	var maxSteps int
	flag.IntVar(&maxSteps, "max-steps", 26501365, "max steps")
	flag.Parse()

	grid := NewGrid(aoc.ReadInputAsByteMatrix())
	start := findStart(grid)
	maxRow, maxCol := grid.Dimensions()

	if maxRow != maxCol {
		log.Fatalf("you have a harder problem than I do")
	}

	entrypointToVisitInfo := map[string]*VisitInfo{}

	for _, rowPos := range []int{0, start.Row, maxRow - 1} {
		for _, colPos := range []int{0, start.Column, maxCol - 1} {
			p := &Position{Row: rowPos, Column: colPos}
			vi := MakeVisitInfo(grid, p, -1)
			entrypointToVisitInfo[p.String()] = vi
		}
	}

	evenVI := entrypointToVisitInfo[start.String()]
	log.Printf("start Reachable (even) = %d", evenVI.Reachable)
	oddVI := entrypointToVisitInfo[Position{start.Row, 0}.String()]
	log.Printf("right Reachable (odd) = %d", oddVI.Reachable)

	log.Printf("down Reachable (odd) = %d", entrypointToVisitInfo[Position{0, start.Column}.String()].Reachable)

	/* Whole grids.  My input has a clear shot north, south, east, and west from
		 the start position, as well as a completely clear border.  This means we
		 can get to any adjacent grid in simple Manhattan distance, so we can
		 ignore obstacles when computing distance between grids.

	         c
		      bXd
		     bXXXd
	      bXXXXXd
		   aXXXSXXXe
	      hXXXXXf
		     hXXXf
		      hXf
		       g

	   Area of complete grids (capital letters) in example is 1+3+5+7+5+3+1= 25

	   25 is 16 + 9; turn this on the side and look the layout like a US flag

	   in general if we can get k grids away in each direction (center to center
	   distance) where we have not enough distance to get to the next center, we
	   can access k**2+(k+1)**2 grids.

	   But for each of the locations 'p', we can reach a fraction of the squares.
	   For the ones on the diagonals, we can reach roughly 1/2 of the squares,
	   and the opposing sides will cancel out.  For the ones in the corners, we
	   can reach 1/4 of the spaces. This isn't quite true if the input has a lot
	   of obstacles, but maybe it's true for us?  In this case, there are an
	   additional 16 things in the example, equal to 2(k+1) total grids.

	   No, that answer is too high, but it's close.

	   However, the remainder when we take our whole grids number, divided by
	   131 (our

	   ... No, that's wrong.  There are two kinds of grids, even and odd.  Start
	   is an even position.  All partial boards are even, because my input allows
	   for 202300 moves before running into a partial board.

	         c
		      bOd
		     bOEOd
	      bOEOEOd
		   aOEOSOEOe
	      hOEOEOf
		     hOEOf
		      hOf
		       g

	*/

	wholePart := maxSteps / maxRow
	remainder := maxSteps % maxRow
	sideGridsPerSide := wholePart - 2
	evenWholeGrids := (wholePart - 2) * (wholePart - 2)
	oddWholeGrids := (wholePart - 1) * (wholePart - 1)
	evenReachable := evenVI.Reachable
	oddReachable := oddVI.Reachable
	wholeReachable := evenWholeGrids*evenReachable + oddWholeGrids*oddReachable

	log.Printf("type a")
	aVI := MakeVisitInfo(grid, &Position{start.Row, maxCol - 1}, remainder-1)
	log.Printf("type c")
	cVI := MakeVisitInfo(grid, &Position{maxRow - 1, start.Column}, remainder-1)
	log.Printf("type e")
	eVI := MakeVisitInfo(grid, &Position{start.Row, 0}, remainder-1)
	log.Printf("type g")
	gVI := MakeVisitInfo(grid, &Position{0, start.Column}, remainder-1)

	// remainder-2 is wrong, to get to the corner we're at a very different number
	bVI := MakeVisitInfo(grid, &Position{maxRow - 1, maxCol - 1}, remainder-2)
	dVI := MakeVisitInfo(grid, &Position{maxRow - 1, 0}, remainder-2)
	fVI := MakeVisitInfo(grid, &Position{0, 0}, remainder-2)
	hVI := MakeVisitInfo(grid, &Position{0, maxCol - 1}, remainder-2)

	sideReachablePerUnit := bVI.Reachable + dVI.Reachable + fVI.Reachable + hVI.Reachable
	sideReachable := sideReachablePerUnit * sideGridsPerSide

	cornerReachable := aVI.Reachable + cVI.Reachable + eVI.Reachable + gVI.Reachable
	log.Printf("cornerReachable = %d\n", cornerReachable)

	totalReachable := cornerReachable + wholeReachable + sideReachable
	fmt.Printf("totalReachable=%d\n", totalReachable)
}
