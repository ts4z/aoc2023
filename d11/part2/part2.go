package main

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick/matrix"
)

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

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

	// this is the way part 1 should have worked anyhow
	// instead of all that row-adding stuff I did because
	// it seemed more straightforward
	const bonus = 999_999 // off by 1 since we still add taxicab distance

	rowBonus := make([]int, m.Rows())
	for i := 0; i < m.Rows(); i++ {
		isRowBlank := true
		for j := 0; isRowBlank && j < m.Columns(); j++ {
			if *m.At(i, j) != '.' {
				isRowBlank = false
			}
		}
		if isRowBlank {
			log.Printf("row blank %d", i)
			rowBonus[i] = bonus
		}
	}

	columnBonus := make([]int, m.Columns())
	for j := 0; j < m.Columns(); j++ {
		isColBlank := true
		for i := 0; isColBlank && i < m.Rows(); i++ {
			if *m.At(i, j) != '.' {
				isColBlank = false
			}
		}
		if isColBlank {
			columnBonus[j] = bonus
			log.Printf("col blank %d", j)
		}
	}

	// for i := 0; i < m.Rows(); i++ {
	// 	for j := 0; j < m.Columns(); j++ {
	// 		fmt.Printf("%c", *m.At(i, j))
	// 	}
	// 	fmt.Printf("\n")
	// }

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

			starI := stars[i]
			starJ := stars[j]

			minRow, maxRow := minMax(starI.Row, starJ.Row)
			for k := minRow + 1; k < maxRow; k++ {
				totalDistance += rowBonus[k]
			}

			minCol, maxCol := minMax(starI.Column, starJ.Column)
			for k := minCol + 1; k < maxCol; k++ {
				totalDistance += columnBonus[k]
			}
		}
	}

	fmt.Printf("%d\n", totalDistance)
}
