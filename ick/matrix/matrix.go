package matrix

// Matrix presents a (potentially) convenient interface for working with an
// array of array of slices allegedly representing a matrix.  It can do things
// like make sure the array is properly rectangular (it doesn't) and invert the
// matrix (it doesn't do that either).  (Maybe someday.)

import (
	"fmt"
	"log"
)

type Matrix[T any] struct {
	a [][]T
}

func (m *Matrix[T]) Rows() int {
	return len(m.a)
}

func (m *Matrix[T]) Columns() int {
	if len(m.a) > 0 {
		return len(m.a[0])
	} else {
		return 0
	}
}

func NewFrom[T any](a [][]T) *Matrix[T] {
	return &Matrix[T]{a: a}
}

func New[T any](rows, cols int) *Matrix[T] {
	m := &Matrix[T]{
		a: make([][]T, rows),
	}
	for i := range m.a {
		m.a[i] = make([]T, cols)
	}
	return m
}

func (m *Matrix[T]) InsertRowWithDefaultValue(before int, def T) {
	nr := make([]T, m.Columns())
	for i := range nr {
		nr[i] = def
	}
	a := m.a
	na := a[0:before]
	na = append(na, nr)
	na = append(na, a[before:]...)
	m.a = na
}

func (m *Matrix[T]) InsertColumnWithDefaultValue(before int, def T) {
	for i := 0; i < m.Rows(); i++ {
		a := m.a[i]
		na := make([]T, len(a)+1)

		copy(na[0:before], a[0:before])
		na[before] = def
		copy(na[1+before:], a[before:])

		m.a[i] = na
	}
}

// Rename, change signature?
func (m *Matrix[T]) At(i, j int) *T {
	return &m.a[i][j]
}

func (m *Matrix[T]) AtAddress(addr Address) *T {
	return &m.a[addr.Row][addr.Column]
}

// Rename, change signature?
func (m *Matrix[T]) Get(i, j int) T {
	return m.a[i][j]
}

func (m *Matrix[T]) GetAddress(addr Address) T {
	return m.a[addr.Row][addr.Column]
}

// Remove.
func (m *Matrix[T]) ForEach(fn func(row, col int)) {
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Columns(); j++ {
			fn(i, j)
		}
	}
}

func (m *Matrix[T]) ForEachAddress(fn func(Address)) {
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Columns(); j++ {
			fn(Address{i, j})
		}
	}
}

type Address struct {
	Row, Column int
}

func posDiff(x, y int) int {
	if x > y {
		return x - y
	} else {
		return y - x
	}
}

func (p Address) TaxicabDistance(there Address) int {
	return posDiff(p.Row, there.Row) + posDiff(p.Column, there.Column)
}

func (p Address) String() string {
	return fmt.Sprintf("Address{%d,%d}", p.Row, p.Column)
}

func (p Address) Go(il RelativeAddress) Address {
	switch il.Direction {
	case Up:
		return Address{p.Row - il.Amount, p.Column}
	case Down:
		return Address{p.Row + il.Amount, p.Column}
	case Left:
		return Address{p.Row, p.Column - il.Amount}
	case Right:
		return Address{p.Row, p.Column + il.Amount}
	default:
		log.Fatalf("bad direction in %#v: %d", il, il.Direction)
		return Address{} // can't happen
	}
}

const (
	Up    = -1
	Down  = 1
	Left  = -2
	Right = 2
)

type Direction int

type RelativeAddress struct {
	Direction Direction
	Amount    int
}

func (m RelativeAddress) Negate() RelativeAddress {
	n := m
	n.Direction = -m.Direction
	return n
}

func (m *Matrix[T]) GetRawMatrix() [][]T {
	return m.a
}
