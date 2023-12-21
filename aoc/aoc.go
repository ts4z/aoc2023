package aoc

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ts4z/aoc2023/argv"
)

func PrintByteMatrix(a [][]byte) {
	PrintByteMatrixTo(os.Stdout, a)
}

func PrintBoolMatrix(a [][]bool) {
	PrintBoolMatrixTo(os.Stdout, a)
}

func PrintByteMatrixTo(w io.Writer, a [][]byte) {
	for _, line := range a {
		for _, ch := range line {
			fmt.Fprintf(w, "%c", ch)
		}
		fmt.Fprintf(w, "\n")
	}
}

func PrintBoolMatrixTo(w io.Writer, a [][]bool) {
	for _, line := range a {
		for _, b := range line {
			if b {
				fmt.Fprintf(w, "X")
			} else {
				fmt.Fprintf(w, ".")
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func ReadInputAsByteMatrix() [][]byte {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	a := make([][]byte, len(lines))
	for i, line := range lines {
		a[i] = []byte(line)
	}
	return a
}

func NewMatrix[T any](rows, cols int) [][]T {
	a := make([][]T, rows)
	for i := range a {
		a[i] = make([]T, cols)
	}
	return a
}

func DimensionsOf(a [][]any) (int, int) {
	return len(a), len(a[0])
}
