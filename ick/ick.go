package ick

// ick: A package for doing things wrong, expediently.
//
// ick: It's like "using namespace std;", but not as polite.
//
// For my AOC2022 projects, I'm sick of pretending this isn't Perl.

import (
	"bufio"
	"errors"
	"io"
	"log" // all kids love log
	"sort"

	"golang.org/x/exp/constraints"
)

type Numeric interface {
	int | int32 | int64 | float64
}

// https://go.dev/doc/tutorial/generics
func Sum[V Numeric](av []V) V {
	var s V
	for _, v := range av {
		s += v
	}
	return s
}

// Sort sorts the data in place.  (This should be renamed QSort?)
func Sort[V constraints.Ordered](av []V) {
	sort.Slice(av, func(i, j int) bool { return av[i] < av[j] })
}

// CSort sorts a copy of the data and returns the updated copy.
// (This should probably be renamed Sort?)
func CSort[V constraints.Ordered](av []V) []V {
	c := DuplicateSlice(av)
	Sort(c)
	return c
}

func DuplicateSlice[V any](av []V) []V {
	cpy := make([]V, len(av))
	copy(cpy, av)
	return cpy
}

func ReadJustOneLine(reader io.Reader) (string, error) {
	r := bufio.NewReader(reader)
	line, err := r.ReadString('\n')
	return line, err
}

func ReadLines(reader io.Reader) ([]string, error) {
	results := []string{}
	r := bufio.NewReader(reader)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return results, nil
		}
		if line[len(line)-1] != '\n' {
			return nil, errors.New("line didn't end in LF")
		}
		line = line[:len(line)-1] // chomp
		results = append(results, line)
	}
}

func ReadBlankLineSeparatedBlocks(reader io.Reader) ([][]string, error) {
	results := [][]string{}
	r := bufio.NewReader(reader)
	done := false
	for !done {
		result := []string{}
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				done = true
				break
			}
			if line[len(line)-1] != '\n' {
				return nil, errors.New("line didn't end in LF")
			}
			line = line[:len(line)-1] // chomp
			if line == "" {
				break
			}
			result = append(result, line)
		}
		results = append(results, result)
	}
	return results, nil
}

// Mustnt is a logging substitute for panic.
func Mustnt(err error) {
	if err != nil {
		log.Fatalf("Mustnt: encountered error: %v", err)
	}
}

// Must extracts the value from a two-argument return and dies if there's an
// error.
func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatalf("Must: encountered error: %v", err)
	}
	return t
}

// MapSlice is the moral equivalent of Lisp's mapcar.  It is not used in the
// presence of polite company by Go programmers, who just write a for loop.
func MapSlice[T, U any](fn func(T) U, ts []T) []U {
	us := make([]U, len(ts))
	for i, t := range ts {
		us[i] = fn(t)
	}
	return us
}

func Keys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func Values[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func New2DArray[T any](n, m int) [][]T {
	a := make([][]T, n)
	for i := 0; i < n; i++ {
		a[i] = make([]T, m)
	}
	return a
}

func New2DArrayWithDefault[T any](n, m int, def T) [][]T {
	a := make([][]T, n)
	for i := 0; i < n; i++ {
		a[i] = make([]T, m)
		for j := 0; j < m; j++ {
			a[i][j] = def
		}
	}
	return a
}
