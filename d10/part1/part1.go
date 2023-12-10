package main

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

type Pipe uint8

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

func charToPipe(ch byte) Pipe {
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

type Grid = Matrix2D[Pipe]
type Viewed = Matrix2D[bool]

func parse() (*Grid, Position) {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	a := ick.New2DArrayWithDefault[Pipe](len(lines), len(lines[0]), Ground)
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

	return &Grid{a: a}, start
}

func process(g *Grid, start Position) int {
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

		log.Printf("at %#v, queue len %d", p, len(queue))

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
	return saw / 2
}

func main() {
	g, start := parse()
	answer := process(g, start)
	fmt.Printf("%d\n", answer)
}
