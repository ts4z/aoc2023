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
	"math/rand"
	"sort"
	"strconv"

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

// Reverse, not in place.
func Reverse[T any](a []T) []T {
	r := make([]T, len(a))
	i := 0
	j := len(a) - 1
	for i < len(a) {
		r[j] = a[i]
		i++
		j--
	}
	log.Printf("input %+v output %+v", a, r)
	return r
}

// Reverse, in place.
func NReverse[T any](a []T) {
	ll := len(a)
	for i, j := 0, ll-1; i < ll/2; i++ {
		a[i], a[j] = a[j], a[i]
		j--
	}
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

// :-(
func Min[T constraints.Ordered](ts []T) T {
	if len(ts) == 0 {
		panic("need at least one value to take a minimum")
	}
	if len(ts) == 1 {
		return ts[0]
	}
	return min(ts[0], Min(ts[1:]))
}

// Like perl's grep.
func Grep[T any](predicate func(T) bool, in []T) []T {
	out := []T{}
	for _, input := range in {
		if predicate(input) {
			out = append(out, input)
		}
	}
	return out
}

func Not[T any](predicate func(T) bool) func(T) bool {
	return func(t T) bool {
		return !predicate(t)
	}
}

func IsEmptyString(s string) bool {
	return len(s) == 0
}

func IsEmptyArray[T any](t []T) bool {
	// now isn't this familiar
	return len(t) == 0
}

func Shuffle[T any](t []T) {
	rand.Shuffle(len(t), func(i, j int) {
		t[i], t[j] = t[j], t[i]
	})
}

func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("can't strconv.Atoi(%q): %v", s, err)
	}
	return n
}
