package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/esc"
	"github.com/ts4z/aoc2023/ick"
)

type PipeCharacter uint8

const (
	Ground = '.' // no connections
	Start  = 'S' // start point in interesting loop
	Dash   = '-' // -, join east and west
	Bar    = '|' // |, join north and south
	Seven  = '7' // 7, join west and south
	Jay    = 'J' // J, join west and north
	Eff    = 'F' // F, join south and east
	Ell    = 'L' // L, join north and east
)

type FloodFillCharacter uint8

const (
	Inside       = '.'
	Outside      = ' '
	PipelinePart = 'X'
)

func charToPipe(ch byte) PipeCharacter {
	switch ch {
	case '.':
		return Ground
	case '-':
		return Dash
	case '|':
		return Bar
	case '7':
		return Seven
	case 'J':
		return Jay
	case 'F':
		return Eff
	case 'L':
		return Ell
	case 'G':
		return Ground
	case 'S':
		return Start
	default:
		log.Fatalf("can't parse %c", ch)
		return Ground // can't happen
	}
}

type Position struct {
	Row, Column int
}

func (p Position) North() Position {
	return Position{p.Row - 1, p.Column}
}

func (p Position) South() Position {
	return Position{p.Row + 1, p.Column}
}

func (p Position) East() Position {
	return Position{p.Row, p.Column + 1}
}

func (p Position) West() Position {
	return Position{p.Row, p.Column - 1}
}

type Matrix2D[T any] struct {
	a [][]T
}

func (m *Matrix2D[T]) Rows() int {
	return len(m.a)
}

func (m *Matrix2D[T]) Columns() int {
	return len(m.a[0])
}

func (m *Matrix2D[T]) Set(p Position, val T) {
	m.a[p.Row][p.Column] = val
}

func (m *Matrix2D[T]) Get(p Position) T {
	return m.a[p.Row][p.Column]
}

func (m *Matrix2D[T]) GetOrZero(p Position) T {
	if p.Row < 0 || p.Row >= m.Rows() ||
		p.Column < 0 || p.Column >= m.Columns() {
		var zero T
		return zero
	}
	return m.a[p.Row][p.Column]
}

func (m *Matrix2D[T]) At(p Position) *T {
	return &m.a[p.Row][p.Column]
}

type PipelineGrid = Matrix2D[PipeCharacter]
type Viewed = Matrix2D[bool]
type FloodFillGrid = Matrix2D[FloodFillCharacter]

func parse() (*PipelineGrid, Position) {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	a := ick.New2DArrayWithDefault[PipeCharacter](len(lines), len(lines[0]), Ground)
	start := Position{-1, -1}

	for i, line := range lines {
		for j, ch := range []byte(line) {
			p := charToPipe(ch)
			a[i][j] = p
			if p == Start {
				start = Position{i, j}
			}
		}
	}

	return &PipelineGrid{a: a}, start
}

func processLoop(g *PipelineGrid, start Position) *Viewed {
	queue := []Position{} // ok it's really gonna be a stack
	saw := 1
	alreadyVisited := &Viewed{a: ick.New2DArray[bool](g.Rows(), g.Columns())}

	*alreadyVisited.At(start) = true

	if p := g.GetOrZero(start.North()); p == Eff || p == Seven || p == Bar {
		log.Printf("initial enqueue North")
		queue = append(queue, start.North())
	}
	if p := g.GetOrZero(start.South()); p == Ell || p == Jay || p == Bar {
		log.Printf("initial enqueue South")
		queue = append(queue, start.South())
	}
	if p := g.GetOrZero(start.East()); p == Jay || p == Seven || p == Dash {
		log.Printf("initial enqueue East")
		queue = append(queue, start.East())
	}
	if p := g.GetOrZero(start.West()); p == Ell || p == Eff || p == Dash {
		log.Printf("initial enqueue West")
		queue = append(queue, start.West())
	}

	if len(queue) != 2 {
		log.Fatalf("suspicious start state with queue len %d", len(queue))
	}

	log.Printf("initial queue %#v", queue)

	for len(queue) != 0 {
		p := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// log.Printf("at %#v, queue len %d", p, len(queue))

		saw++
		*alreadyVisited.At(p) = true
		ch := *g.At(p)
		nexts := []Position{}
		switch ch {
		case Seven:
			nexts = append(nexts, p.West(), p.South())
		case Eff:
			nexts = append(nexts, p.South(), p.East())
		case Ell:
			nexts = append(nexts, p.North(), p.East())
		case Jay:
			nexts = append(nexts, p.North(), p.West())
		case Bar:
			nexts = append(nexts, p.North(), p.South())
		case Dash:
			nexts = append(nexts, p.East(), p.West())
		case Start:
			log.Printf("not returning to start")
		default:
			log.Fatalf("at position %#v, found ch %d which is not part of the pipeline",
				p, ch)
		}

		for _, next := range nexts {
			if alreadyVisited.Get(next) {
				continue
			}
			queue = append(queue, next)
		}
	}

	if saw&1 != 0 {
		// this check is wrong, I don't know why
		log.Printf("saw is odd, which can't happen")
	}

	return alreadyVisited
}

func maskLoop(g *PipelineGrid, v *Viewed) *PipelineGrid {
	r := &PipelineGrid{ick.New2DArrayWithDefault(g.Rows(), g.Columns(), PipeCharacter(Ground))}
	for i := 0; i < g.Rows(); i++ {
		for j := 0; j < g.Columns(); j++ {
			at := Position{i, j}
			if v.Get(at) {
				*r.At(at) = g.Get(at)
			}
		}
	}
	return r
}

func expandGrid(in *PipelineGrid) *PipelineGrid {
	out := &PipelineGrid{ick.New2DArrayWithDefault(in.Rows()*2+1, in.Columns()*2+1, PipeCharacter(Ground))}
	for i := 1; i < out.Rows()-1; i++ {
		for j := 1; j < out.Columns()-1; j++ {
			evenRow := i&1 == 0
			evenCol := j&1 == 0
			oddRow := !evenRow
			oddCol := !evenCol
			here := Position{i, j}

			// log.Printf("here %#v", here)

			if oddRow && oddCol {
				*out.At(here) = in.Get(Position{i / 2, j / 2})
				continue
			}

			if evenRow && evenCol {
				continue // just space here
			}

			// we have sorted data to our north and west; we can
			// look at our output in those locations.

			if evenRow /* and odd column */ {
				north := out.Get(here.North())
				if north == Bar || north == Seven || north == Eff {
					*out.At(here) = '|'
				}
				continue
			}

			if evenCol /* and odd row */ {
				west := out.Get(here.West())
				if west == Dash || west == Eff || west == Ell {
					*out.At(here) = '-'
				}
				continue
			}

			log.Fatalf("can't get here")
		}
	}
	return out
}

func ToFloodFillMap(v *PipelineGrid) *FloodFillGrid {
	a := ick.New2DArrayWithDefault[FloodFillCharacter](v.Rows(), v.Columns(), Inside)
	r := &FloodFillGrid{a: a}
	for i := 0; i < r.Rows(); i++ {
		for j := 0; j < r.Columns(); j++ {
			at := Position{i, j}
			if v.Get(at) == '.' {
				*r.At(at) = Inside
			} else {
				*r.At(at) = PipelinePart
			}
		}
	}
	return r
}

func contractGrid(in *FloodFillGrid) *FloodFillGrid {
	a := ick.New2DArrayWithDefault[FloodFillCharacter](in.Rows()/2, in.Columns()/2, Outside)
	out := &FloodFillGrid{a: a}
	for i := 0; i < out.Rows(); i++ {
		for j := 0; j < out.Columns(); j++ {
			*out.At(Position{i, j}) = in.Get(Position{2*i + 1, 2*j + 1})
		}
	}
	return out
}

type floodFillUpdate struct {
	p         Position
	queueSize int
}

func ProcessFloodFill(ff *FloodFillGrid) {
	updates := make(chan floodFillUpdate, 10000)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		printFloodFillMapTripleWide(ff, updates)
		wg.Done()
	}()

	// init
	outside := []Position{}
	noItIsnt := func(p Position) {
		if ff.Get(p) == Inside {
			// no it isn't
			*ff.At(p) = Outside
			outside = append(outside, p)
			updates <- floodFillUpdate{p: p, queueSize: len(outside)}
		}
	}

	// prime pump, starting from outside
	rows := ff.Rows()
	columns := ff.Columns()

	for i := 0; i < columns; i++ {
		outside = append(outside, Position{Row: 0, Column: i})
		outside = append(outside, Position{Row: rows - 1, Column: i})
	}

	for j := 0; j < columns-1; j++ {
		outside = append(outside, Position{Row: j, Column: 0})
		outside = append(outside, Position{Row: j, Column: columns - 1})
	}

	for len(outside) != 0 {
		// Maybe don't do it in the obvious order?
		// (this makes it more fun to watch)
		if true {
			k := rand.Intn(len(outside))
			n := len(outside) - 1
			outside[k], outside[n] = outside[n], outside[k]
		}

		p := outside[len(outside)-1]
		outside = outside[:len(outside)-1]
		if ff.GetOrZero(p.North()) == Inside {
			noItIsnt(p.North())
		}
		if ff.GetOrZero(p.South()) == Inside {
			noItIsnt(p.South())
		}
		if ff.GetOrZero(p.East()) == Inside {
			noItIsnt(p.East())
		}
		if ff.GetOrZero(p.West()) == Inside {
			noItIsnt(p.West())
		}
	}

	// wait so all the updates are printed and we don't mangle the output
	close(updates)
	wg.Wait()
}

func printMap[T ~uint8](ff *Matrix2D[T]) {
	fmt.Printf("inside/outside map:\n")
	for i := 0; i < ff.Rows(); i++ {
		for j := 0; j < ff.Columns(); j++ {
			at := Position{i, j}
			ch := ff.Get(at)
			fmt.Printf("%c", rune(ch))
		}
		fmt.Printf("\n")
	}
}

func printFloodFillMap(ff *FloodFillGrid) int {
	inside := 0
	outside := 0
	pipeline := 0
	fmt.Printf("inside/outside map:\n")
	for i := 0; i < ff.Rows(); i++ {
		for j := 0; j < ff.Columns(); j++ {
			at := Position{i, j}
			ch := ff.Get(at)
			switch ch {
			case Inside:
				inside++
			case Outside:
				outside++
			case PipelinePart:
				pipeline++
			default:
				log.Fatalf("wtf at %#v ch=%c", at, ch)
			}
			fmt.Printf("%c", ch)
		}
		fmt.Printf("\n")
	}

	log.Printf("inside %d outside %d pipeline %d",
		inside, outside, pipeline)

	return inside
}

func printFloodFillMapTripleWide(ff *FloodFillGrid, outsides chan floodFillUpdate) {
	fmt.Printf("inside/outside map:\n")

	esc.Clear()
	esc.Home()

	inside := 0
	outside := 0

	printFull := func() {
		for i := 0; i < ff.Rows(); i++ {
			for j := 0; j < ff.Columns(); j++ {
				at := Position{i, j}
				ch := ff.Get(at)
				if ch == Inside {
					inside++
				}
				fmt.Printf("%c%c%c", ch, ch, ch)
			}
			fmt.Printf("\n")
		}
	}

	printFull()

	for up := range outsides {
		time.Sleep(1 * time.Millisecond)
		esc.Home()
		inside--
		outside++
		fmt.Printf("%d inside %d outside, %d in queue          ", inside, outside, up.queueSize)
		p := up.p
		esc.Goto(p.Row+1, 3*p.Column+1)
		fmt.Printf("   ")
	}

	esc.Home()
	printFull()
}

// fixStart replaces the start position with the pipe that it represents.
func fixStart(g *PipelineGrid, start Position) {
	north := g.Get(start.North())
	south := g.Get(start.South())
	east := g.Get(start.East())
	west := g.Get(start.West())
	northCxn := north == Seven || north == Eff || north == Bar
	southCxn := south == Ell || south == Jay || south == Bar
	eastCxn := east == Seven || east == Jay || east == Dash
	westCxn := west == Eff || west == Ell || west == Dash
	found := 0
	if northCxn && southCxn {
		*g.At(start) = Bar
		found++
	}
	if northCxn && eastCxn {
		*g.At(start) = Ell
		found++
	}
	if northCxn && westCxn {
		*g.At(start) = Jay
		found++
	}
	if southCxn && eastCxn {
		*g.At(start) = Eff
		found++
	}
	if southCxn && westCxn {
		*g.At(start) = Seven
		found++
	}
	if eastCxn && westCxn {
		*g.At(start) = Dash
		found++
	}
	if found != 1 {
		log.Fatalf("found %d, expected 1, there's a bug", found)
	}
}

func main() {
	g, start := parse()

	// recover pipeline shape at start
	fixStart(g, start)

	viewed := processLoop(g, start)

	fmt.Printf("visited by pipeline:\n")
	for i := 0; i < viewed.Rows(); i++ {
		for j := 0; j < viewed.Columns(); j++ {
			if viewed.Get(Position{i, j}) {
				fmt.Printf("X")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}

	masked := maskLoop(g, viewed)
	fmt.Printf("masked pipeline:\n")
	printMap(masked)

	expanded := expandGrid(masked)
	fmt.Printf("expanded masked pipeline:\n")
	printMap(expanded)

	ff := ToFloodFillMap(expanded)

	fmt.Printf("initial floodfill map:\n")
	printFloodFillMap(ff)

	ProcessFloodFill(ff)

	fmt.Printf("final floodfill map:\n")
	printFloodFillMap(ff)

	contracted := contractGrid(ff)

	inside := printFloodFillMap(contracted)

	fmt.Printf("answer %d\n", inside)
}
