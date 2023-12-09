package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/lpc"
)

func parse() ([][]int, error) {
	var r [][]int
	var firstErr error
	for lpc := lpc.New(ick.Must(argv.ReadChompAll())); !lpc.EOF(); lpc.Next() {
		line := lpc.Current()
		parse := func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil && firstErr == nil {
				firstErr = lpc.Wrap("converting string to int", err)
			}
			return i
		}
		a := ick.MapSlice(parse, strings.Split(line, " "))
		r = append(r, a)
	}
	return r, firstErr
}

func computeDiffs(input []int) []int {
	diffs := []int{}
	for i := 1; i < len(input); i++ {
		diffs = append(diffs, input[i]-input[i-1])
	}
	return diffs
}

func isAllZero(a []int) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}

func processInput(a []int) int {
	diffs := computeDiffs(a)
	r := 0
	if !isAllZero(a) {
		aa := processInput(diffs)
		r = a[len(a)-1] + aa
	}
	log.Printf("%+v %d", a, r)
	return r
}

func processInputs(inputs [][]int) int {
	r := 0
	for _, input := range inputs {
		log.Printf("[begin] %+v ---\n", input)
		v := processInput(input)
		log.Printf("[-end-] %+v next %d\n", input, v)
		r += v
	}
	return r
}

func main() {
	inputs, err := parse()
	if err != nil {
		log.Fatalf("can't parse input: %v", err)
	}
	answer := processInputs(inputs)
	fmt.Printf("%d\n", answer)
}
