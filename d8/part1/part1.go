package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/lpc"
)

var lineRE = regexp.MustCompile(`^(...) = \((...), (...)\)$`)

type Node struct {
	Left, Right string
}

type Input struct {
	Directions string
	Nodes      map[string]*Node
}

func parse() (*Input, error) {
	lpc := lpc.New(ick.Must(argv.ReadChompAll()))
	nodes := map[string]*Node{}
	directions := lpc.Current()
	lpc.Next()
	lpc.MustEatBlankLine()
	for ; !lpc.EOF(); lpc.Next() {
		matches := lineRE.FindSubmatch([]byte(lpc.Current()))
		if matches == nil {
			return nil, lpc.Wrap("reading node line", fmt.Errorf("RE doesn't match %q", lpc.Current()))
		}
		nodes[string(matches[1])] = &Node{Left: string(matches[2]), Right: string(matches[3])}
	}

	r := &Input{directions, nodes}
	return r, nil
}

func main() {
	input, err := parse()
	if err != nil {
		log.Fatalf("can't parse: %v", err)
	}

	at := "AAA"
	step := 0
	for ; at != "ZZZ"; step++ {
		mod := step % len(input.Directions)
		thisStep := input.Directions[mod]
		n := input.Nodes[at]
		if n == nil {
			log.Fatalf("step %d (%d) can't find node %s", step, mod, at)
		}
		var next string
		if thisStep == 'L' {
			next = n.Left
		} else if thisStep == 'R' {
			next = n.Right
		} else {
			log.Fatalf("step %d (%d) is %v?", step, mod, thisStep)
			next = ""
		}

		log.Printf("step %d (%d) steps %c to %v", step, mod, thisStep, next)
		at = next
	}

	log.Printf("%d steps", step)
}
