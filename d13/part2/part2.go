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

type reflectionInfo struct {
	reflectionAt     int
	reflectionBegins int
	reflectionEnds   int
}

func findReflection(lines []int, what string) ([]reflectionInfo, bool) {
	answers := []reflectionInfo{}
	for i := 0; i < len(lines)-1; i++ {
		if lines[i] == lines[i+1] {
			log.Printf("two in a row at %s %d\n", what, i)
			j := i
			k := i + 1
			for {
				j--
				k++
				if j < 0 || k >= len(lines) {
					answers = append(answers, reflectionInfo{i, j + 1, k - 1})
					break
				}
				if lines[j] != lines[k] {
					log.Printf("j=%d k=%d mismatch", j, k)
					break
				}
			}
		}
	}
	return answers, len(answers) != 0
}

func flip(ch byte) byte {
	if ch == '#' {
		return '.'
	} else if ch == '.' {
		return '#'
	} else {
		log.Fatalf("can't flip character %c", ch)
		return '?'
	}
}

func printMatrix(m *matrix.Matrix[byte]) {
	// for i := 0; i < m.Rows(); i++ {
	// 	for j := 0; j < m.Columns(); j++ {
	// 		fmt.Printf("%c", *m.At(i, j))
	// 	}
	// 	fmt.Printf("\n")
	// }
}

func processAllPossibleSmudges(n int, lines []string) int {
	m := toMatrix(lines)

	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Columns(); j++ {
			was := *m.At(i, j)
			*m.At(i, j) = flip(was)
			fmt.Printf("flipped %d,%d:\n", i, j)
			printMatrix(m)
			fmt.Printf("\n")
			score, ok := processMatrix(m, i, j)
			*m.At(i, j) = was
			if ok {
				fmt.Printf("restored %d,%d:\n", i, j)
				printMatrix(m)
				log.Printf("flipped %d,%d, found result %d", i, j, score)
				return score
			}
		}
	}

	log.Printf("can't find reflection in matrix %d", n)
	log.Printf("=> .... 012345678901234567890123456789")
	for i, line := range lines {
		log.Printf("=> [%2d] %s", i, line)
	}
	log.Printf("=> .... 012345678901234567890123456789")
	log.Fatalf("can't find reflection")
	return -1 // can't happen
}

func processMatrix(m *matrix.Matrix[byte], mutatedRow, mutatedColumn int) (int, bool) {
	fmt.Printf("input matrix: \n")
	printMatrix(m)
	fmt.Printf("\n")

	{
		rows := []int{}
		for i := 0; i < m.Rows(); i++ {
			k := 0
			for j := 0; j < m.Columns(); j++ {
				maybeOne := 0
				kWas := k
				if *m.At(i, j) == '#' {
					maybeOne = 1
				}
				k = (k << 1) | maybeOne
				if k < kWas {
					log.Fatalf("wrap")
				}
			}
			rows = append(rows, k)
		}

		log.Printf("rows %+v", rows)

		if ris, ok := findReflection(rows, "row"); ok {
			for _, ri := range ris {
				if mutatedRow < 0 || mutatedRow >= ri.reflectionBegins && mutatedRow < ri.reflectionEnds {
					return (1 + ri.reflectionAt) * 100, true
				} else {
					log.Printf("row reflection outside area %+v", ri)
				}
			}
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
			kWas := k
			k = (k << 1) | maybeOne
			if k < kWas {
				log.Fatalf("wrap")
			}
		}
		columns = append(columns, k)
	}

	log.Printf("columns %+v", columns)

	if ris, ok := findReflection(columns, "column"); ok {
		for _, ri := range ris {
			if mutatedColumn < 0 || mutatedColumn >= ri.reflectionBegins && mutatedColumn < ri.reflectionEnds {
				return 1 + ri.reflectionAt, true
			} else {
				log.Printf("column reflection outside area %+v", ri)
			}
		}
	}

	// log.Printf("can't find reflection in this matrix")
	// log.Printf("=> .... 012345678901234567890123456789")
	// for i, line := range lines {
	// 	log.Printf("=> [%2d] %s", i, line)
	// }
	// log.Printf("=> .... 012345678901234567890123456789")
	// log.Printf("can't find reflection")

	return -1, false
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

	total := 0
	for i, block := range blocks {
		total += processAllPossibleSmudges(i, block)
	}

	fmt.Printf("total %d\n", total)
}
