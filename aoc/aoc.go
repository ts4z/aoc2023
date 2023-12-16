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

func PrintByteMatrixTo(w io.Writer, a [][]byte) {
	for _, line := range a {
		for _, ch := range line {
			fmt.Fprintf(w, "%c", ch)
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
