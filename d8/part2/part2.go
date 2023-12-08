/*

I took a spoiler here having worked on this for a couple hours.  Once I found
the length of the cycles, I took the LCM of that number and it happened to be
correct.

I don't understand why the LCM of the cycles is significant, however,
the number of steps before a cycle is detected is very small, and I didn't
look at the number of cycles to get to Z nodes.

I found this a very frustrating, disappointing answer.

*/

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

// I once knew a dog named 3B.
type ThreeB [3]byte

func (tb ThreeB) String() string {
	return string(tb[:])
}

type Node struct {
	Left, Right ThreeB
}

type Input struct {
	Directions string
	Nodes      map[ThreeB]Node
}

func to3Byte(b []byte) ThreeB {
	if len(b) != 3 {
		log.Fatalf("input %v has invalid length", b)
	}
	return ThreeB{b[0], b[1], b[2]}
}

func parse() (*Input, error) {
	lpc := lpc.New(ick.Must(argv.ReadChompAll()))
	nodes := map[ThreeB]Node{}
	directions := lpc.Current()
	lpc.Next()
	lpc.MustEatBlankLine()
	for ; !lpc.EOF(); lpc.Next() {
		matches := lineRE.FindSubmatch([]byte(lpc.Current()))
		if matches == nil {
			return nil, lpc.Wrap("reading node line", fmt.Errorf("RE doesn't match %q", lpc.Current()))
		}
		nodes[to3Byte(matches[1])] = Node{
			Left:  to3Byte(matches[2]),
			Right: to3Byte(matches[3]),
		}
	}

	r := &Input{directions, nodes}
	return r, nil
}

func allEndInZ(a []ThreeB) bool {
	for _, s := range a {
		if s[2] != 'Z' {
			return false
		}
	}
	return true
}

type Position struct {
	Mod      int
	NodeName ThreeB
}

type CycleInfo struct {
	Reported      bool
	FirstStepHere int
}

func main() {
	input, err := parse()
	if err != nil {
		log.Fatalf("can't parse: %v", err)
	}

	ats := ick.Grep(func(s ThreeB) bool { return s[2] == 'A' }, ick.Keys(input.Nodes))

	step := 0
	mod := 0

	cycleMap := map[Position]CycleInfo{}
	cycleCounts := map[int]int{}

	for !allEndInZ(ats) {

		if (step & 0xFFFFF) == 0xFFFFF {
			log.Printf("step %d %+v", step, ats)
		}

		if mod >= len(input.Directions) {
			mod = 0
		}
		for i, at := range ats {

			if cycleInfo, ok := cycleMap[Position{Mod: mod, NodeName: at}]; ok {
				if !cycleInfo.Reported {
					length := step - cycleInfo.FirstStepHere
					// log.Printf("position %v,%v cycles %v,%v length %v", mod, at,
					// 	cycleInfo.FirstStepHere, step, length)
					cycleInfo.Reported = true
					cycleMap[Position{Mod: mod, NodeName: at}] = cycleInfo

					if _, ok := cycleCounts[length]; !ok {
						cycleCounts[length] = step - length
						log.Printf("%d cycleCounts=%#v", i, cycleCounts)
					}
				}
			} else {
				cycleMap[Position{Mod: mod, NodeName: at}] = CycleInfo{FirstStepHere: step}
			}

			thisStep := input.Directions[mod]

			if thisStep == 'L' {
				ats[i] = input.Nodes[at].Left
			} else if thisStep == 'R' {
				ats[i] = input.Nodes[at].Right
			} else {
				log.Fatalf("step %d (%d) is %v?", step, mod, thisStep)
			}
		}

		mod++
		step++
	}

	log.Printf("%d steps and we're at %+v", step, ick.MapSlice(func(tb ThreeB) string { return string(tb[:]) }, ats))
}
