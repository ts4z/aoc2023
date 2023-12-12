package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

func processConsumingSeq(s string, seq []int) int {
	fmt.Printf("   processConsumingSeq(%q, %+v)\n", s, seq)
	if len(s) < seq[0] {
		return 0
	}

	span := seq[0]

	for at := 0; at < span; at++ {
		if s[at] == '#' || s[at] == '?' {
			// good
		} else {
			fmt.Printf("    no span")
			return 0
		}
	}

	seq = seq[1:]
	s = s[span:]

	if len(seq) == 0 && len(s) == 0 {
		return 1
	}
	if len(s) == 0 {
		return 0
	}

	if s[0] == '#' {
		// Sequence must be terminated by . or ?, not spill to another #
		fmt.Printf("    looking at # after seq, no good")
		return 0
	}

	// chop off the character we just looked at
	s = s[1:]

	// looks good so far, continue
	fmt.Printf("    continuing")
	return process(s, seq)
}

func process(s string, seq []int) int {
	fmt.Printf("   process(%q, %+v)\n", s, seq)
	if len(s) == 0 {
		if len(seq) == 0 {
			fmt.Printf("   process(%q, %+v) => %d #1\n", s, seq, 1)
			return 1
		} else {
			fmt.Printf("   process(%q, %+v) => %d #2\n", s, seq, 0)
			return 0
		}
	}

	if len(seq) == 0 {
		if strings.Index(s, "#") >= 0 {
			fmt.Printf("   process(%q, %+v) => %d #3\n", s, seq, 0)
			return 0
		} else {
			fmt.Printf("   process(%q, %+v) => %d #4\n", s, seq, 1)
			return 1
		}
	}

	if s[0] == '.' {
		fmt.Printf("    recurse to process\n")
		return process(s[1:], seq)
	}

	if s[0] == '#' {
		fmt.Printf("    recurse to consuming\n")
		return processConsumingSeq(s, seq)
	}

	if s[0] == '?' {
		fmt.Printf("      two-way recurse(%q, %+v)\n", s, seq)
		no := process(s[1:], seq)
		yes := processConsumingSeq(s, seq)
		fmt.Printf("      two-way recurse(%q, %+v) returned %d+%d\n", s, seq, no, yes)
		return no + yes
	}

	log.Fatalf("can't happen: invalid pic character in %q", s)
	return 0
}

func main() {
	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read input: %v", err)
	}

	total := 0
	for _, line := range lines {
		parts := strings.Split(line, " ")
		seqs := ick.MapSlice(func(n string) int {
			return ick.Must(strconv.Atoi(n))
		}, strings.Split(parts[1], ","))
		sub := process(parts[0], seqs)
		fmt.Printf("process %q %+v => %d\n", parts[0], seqs, sub)
		total += sub
	}
	fmt.Printf("total %d\n", total)
}
