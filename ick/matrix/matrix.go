package matrix

type Matrix[T any] struct {
	a [][]T
}

func (m *Matrix[T]) Rows() int {
	return len(m.a)
}

func (m *Matrix[T]) Columns() int {
	return len(m.a[0])
}

func NewFrom[T any](a [][]T) *Matrix[T] {
	return &Matrix[T]{a: a}
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

func (m *Matrix[T]) At(i, j int) *T {
	return &m.a[i][j]
}

func (m *Matrix[T]) Get(i, j int) T {
	return m.a[i][j]
}

func (m *Matrix[T]) ForEach(fn func(row, col int)) {
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Columns(); j++ {
			fn(i, j)
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
