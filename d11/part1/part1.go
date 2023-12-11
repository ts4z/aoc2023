package main

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/ick/matrix"
)

func main() {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	rows := [][]byte{}
	for _, line := range lines {
		rows = append(rows, ([]byte(line)))
	}

	m := matrix.NewFrom(rows)

	blankRows := []int{}
	for i := 0; i < m.Rows(); i++ {
		isRowBlank := true
		for j := 0; isRowBlank && j < m.Columns(); j++ {
			if *m.At(i, j) != '.' {
				isRowBlank = false
			}
		}
		if isRowBlank {
			log.Printf("row blank %d", i)
			blankRows = append(blankRows, i)
		}
	}

	blankCols := []int{}
	for j := 0; j < m.Columns(); j++ {
		isColBlank := true
		for i := 0; isColBlank && i < m.Rows(); i++ {
			if *m.At(i, j) != '.' {
				isColBlank = false
			}
		}
		if isColBlank {
			log.Printf("col blank %d", j)
			blankCols = append(blankCols, j)
		}
	}

	ick.NReverse(blankRows)
	for _, row := range blankRows {
		m.InsertRowWithDefaultValue(row, '.')
	}

	ick.NReverse(blankCols)
	for _, col := range blankCols {
		m.InsertColumnWithDefaultValue(col, '.')
	}

	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Columns(); j++ {
			fmt.Printf("%c", *m.At(i, j))
		}
		fmt.Printf("\n")
	}

	stars := []matrix.Address{}
	m.ForEach(func(i, j int) {
		if *m.At(i, j) == '#' {
			stars = append(stars, matrix.Address{i, j})
		}
	})

	totalDistance := 0
	for i := 0; i < len(stars)-1; i++ {
		for j := i + 1; j < len(stars); j++ {
			totalDistance += stars[i].TaxicabDistance(stars[j])
		}
	}

	fmt.Printf("%d\n", totalDistance)
}
