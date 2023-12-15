package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ts4z/aoc2023/argv"
)

func hash(s string) int {
	h := 0
	for _, ch := range []byte(s) {
		h += int(ch)
		h *= 17
		h %= 256
	}
	return h
}

func main() {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	total := 0
	insts := strings.Split(lines[0], ",")

	for _, inst := range insts {
		total += hash(inst)
	}

	fmt.Printf("%d\n", total)
}
