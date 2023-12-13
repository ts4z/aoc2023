package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/ick/matrix"
)

func toMatrix(lines []string) *matrix.Matrix[byte] {
	rows := [][]byte{}
	for _, line := range lines {
		rows = append(rows, ([]byte(line)))
	}

	return matrix.NewFrom(rows)
}

func findReflection(lines []int) (int, bool) {
	for i := 0; i < len(lines)-1; i++ {
		if lines[i] == lines[i+1] {
			log.Printf("two in a row at %d\n", i)
			j := i - 1
			k := i + 2
			for {
				if j < 0 || k > len(lines)-1 {
					return i + 1, true
				}
				if lines[j] != lines[k] {
					log.Printf("j=%d k=%d mismatch", j, k)
					// no reflection
					break
				}
				j--
				k++
			}
		}
	}
	return -1, false
}

func processBlock(lines []string) int {
	m := toMatrix(lines)
	{
		rows := []int{}
		for i := 0; i < m.Rows(); i++ {
			k := 0
			for j := 0; j < m.Columns(); j++ {
				maybeOne := 0
				if *m.At(i, j) == '#' {
					maybeOne = 1
				}
				k = (k << 1) | maybeOne
			}
			rows = append(rows, k)
		}

		log.Printf("rows %+v", rows)

		if at, ok := findReflection(rows); ok {
			return 100 * at
		}
	}

	columns := []int{}
	for j := 0; j < m.Columns(); j++ {
		k := 0
		for i := 0; i < m.Rows(); i++ {
			maybeOne := 0
			if *m.At(i, j) == '#' {
				maybeOne = 1
			}
			k = (k << 1) | maybeOne
		}
		columns = append(columns, k)
	}

	log.Printf("columns %+v", columns)

	if at, ok := findReflection(columns); ok {
		return at
	}

	log.Printf("can't find reflection in this matrix")
	log.Printf("=> .... 012345678901234567890123456789")
	for i, line := range lines {
		log.Printf("=> [%2d] %s", i, line)
	}
	log.Printf("=> .... 012345678901234567890123456789")
	log.Printf("can't find reflection")

	return 0
}

func readInput() [][]string {
	blocks, err := ick.ReadBlankLineSeparatedBlocks(os.Stdin)
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}
	return blocks
}

func main() {
	// this should read ARGV, perl style, but I didn't implement that yet
	blocks := readInput()

	total := ick.Sum(ick.MapSlice(processBlock, blocks))

	fmt.Printf("total %d\n", total)
}
