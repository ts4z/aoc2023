package main

// The initial implementation from part1 was pretty slow, so I added some
// threads.  This was still slow, so I added memoization.
//
// It finished almost immediately.
//
// Of course memoizing makes it faster.  It prunes all the hard states at the
// end of the tree that we're calling recursively.
//
// No attempt is made to prevent write-write cache stomping.  No attempt is
// made to constrain the size of the cache.
//
// I changed the code down to 1 thread and that's faster.  (This of course
// fixes issues with conflicts.)

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

func processConsumingSeq(s string, seq []int) int {
	// fmt.Printf("   processConsumingSeq(%q, %+v)\n", s, seq)
	if len(s) < seq[0] {
		return 0
	}

	span := seq[0]

	for at := 0; at < span; at++ {
		if s[at] == '#' || s[at] == '?' {
			// good
		} else {
			// fmt.Printf("    no span")
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
		// fmt.Printf("    looking at # after seq, no good")
		return 0
	}

	// chop off the character we just looked at
	s = s[1:]

	// looks good so far, continue
	// fmt.Printf("    continuing")
	return process(s, seq)
}

var processCache sync.Map

func process(s string, seq []int) int {
	key := fmt.Sprintf("%s %+v", s, seq)
	if value, ok := processCache.Load(key); ok {
		return value.(int)
	} else {
		answer := processWithoutCaching(s, seq)
		processCache.Store(key, answer)
		return answer
	}
}

func processWithoutCaching(s string, seq []int) int {
	// fmt.Printf("   process(%q, %+v)\n", s, seq)
	if len(s) == 0 {
		if len(seq) == 0 {
			// fmt.Printf("   process(%q, %+v) => %d #1\n", s, seq, 1)
			return 1
		} else {
			// fmt.Printf("   process(%q, %+v) => %d #2\n", s, seq, 0)
			return 0
		}
	}

	if len(seq) == 0 {
		if strings.Index(s, "#") >= 0 {
			// fmt.Printf("   process(%q, %+v) => %d #3\n", s, seq, 0)
			return 0
		} else {
			// fmt.Printf("   process(%q, %+v) => %d #4\n", s, seq, 1)
			return 1
		}
	}

	if s[0] == '.' {
		// fmt.Printf("    recurse to process\n")
		return process(s[1:], seq)
	}

	if s[0] == '#' {
		// fmt.Printf("    recurse to consuming\n")
		return processConsumingSeq(s, seq)
	}

	if s[0] == '?' {
		// fmt.Printf("      two-way recurse(%q, %+v)\n", s, seq)
		no := process(s[1:], seq)
		yes := processConsumingSeq(s, seq)
		// fmt.Printf("      two-way recurse(%q, %+v) returned %d+%d\n", s, seq, no, yes)
		return no + yes
	}

	log.Fatalf("can't happen: invalid pic character in %q", s)
	return 0
}

func main() {
	lines := make(chan string, 100)
	errs := make(chan error, 1)

	argv.ReadToChannels(lines, errs)

	go func() {
		err := <-errs
		log.Printf("error on input: %v", err)
	}()

	type lineAndNumber struct {
		line   string
		number int
	}

	lineChannel := make(chan lineAndNumber, 100)

	go func() {
		i := 0
		for line := range lines {
			lineChannel <- lineAndNumber{line, i}
			i++
		}
		close(lineChannel)
	}()

	subtotalChannel := make(chan int, 100)

	go func() {
		var wg sync.WaitGroup
		// Spawn any number of processes here; however, 1 seems to work best.
		for i := 0; i < 1; i++ {
			wg.Add(1)
			go func() {
				for input := range lineChannel {
					line := input.line
					parts := strings.Split(line, " ")
					seqs := ick.MapSlice(func(n string) int {
						return ick.Atoi(n)
					}, strings.Split(parts[1], ","))
					s5 := parts[0] + "?" + parts[0] + "?" + parts[0] + "?" + parts[0] + "?" + parts[0]
					seq5 := []int{}
					seq5 = append(seq5, seqs...)
					seq5 = append(seq5, seqs...)
					seq5 = append(seq5, seqs...)
					seq5 = append(seq5, seqs...)
					seq5 = append(seq5, seqs...)
					sub := process(s5, seq5)
					fmt.Printf("(%d) process %q %+v => %d\n", input.number,
						parts[0], seqs, sub)
					subtotalChannel <- sub
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(subtotalChannel)
	}()

	total := 0
	for sub := range subtotalChannel {
		total += sub
	}
	fmt.Printf("total %d\n", total)

	entries := 0
	keyLen := 0
	processCache.Range(func(k any, v any) bool {
		keyLen += len(k.(string))
		entries++
		return true
	})
	fmt.Printf("%d entries in cache, keys consuming %d bytes of data\n", entries, keyLen)
}
