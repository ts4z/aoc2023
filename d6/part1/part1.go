package main

import (
	"log"
	"regexp"
	"strconv"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

func mustAtoi(s string) int {
	return ick.Must(strconv.Atoi(s))
}

var spaceRE = regexp.MustCompile("\\s+")

func parseNumbers(s string) []int {
	return ick.MapSlice(mustAtoi, ick.Grep(ick.Not(ick.IsEmptyString), spaceRE.Split(s, -1)))
}

func main() {
	lines := ick.Must(argv.ReadChompAll())
	type reader struct {
		prefix   string
		consumer func(string) error
	}
	var times []int
	var distances []int
	readers := []reader{
		{"Time:", func(line string) error {
			times = parseNumbers(line)
			return nil
		}},
		{"Distance:", func(line string) error {
			distances = parseNumbers(line)
			return nil
		}},
	}

	for _, line := range lines {
		found := false
		for _, r := range readers {
			log.Printf("%q %q", line[:len(r.prefix)], r.prefix)
			if line[:len(r.prefix)] == r.prefix {
				found = true
				err := r.consumer(line[len(r.prefix):])
				if err != nil {
					log.Fatalf("fail consuming %q by %v: %v", line, r, err)
				}
				break
			}
		}
		if !found {
			log.Fatalf("no reader for line: %q", line)
		}
	}

	log.Printf("times %+v", times)
	log.Printf("dists %+v", times)

	if len(times) != len(distances) {
		log.Fatalf("length mismatch")
	}

	victoryProduct := 1
	for i := 0; i < len(times); i++ {
		t := times[i]
		d := distances[i]
		victories := 0
		for j := 1; j < t-1; j++ {
			travelled := (t - j) * j
			if travelled > d {
				victories++
			}
		}
		victoryProduct *= victories
	}
	log.Printf("victoryProduct: %d", victoryProduct)
}
