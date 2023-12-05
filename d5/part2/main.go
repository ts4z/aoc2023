package main

/* aoc 2023 day 5 part 2

   This works but it's slow (~10 minutes).  It can probably be trivially sped
   up with some hacks, but it's brute force.  It would be better to track
   batches of things and pass slices around, but that means splitting when we
   find a partial match of a range, and I lost interest because the program
   finished running and I had my answer.
*/

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/lpc"
)

var doShuffle = false

type MapLine struct {
	InputStart  int
	OutputStart int
	Length      int
}

type Map struct {
	Desc     string
	From     string
	To       string
	MapLines []MapLine
}

type InputFile struct {
	Seeds []int
	Maps  map[string]*Map
}

// asscanf is fmt.Sscanf without short reads ("all sscanf").

func asscanf(s string, fm string, args ...any) error {
	n, err := fmt.Sscanf(s, fm, args...)
	if err != nil {
		return err
	}
	expected := 0
	for i := 0; i < len(fm); i++ {
		if fm[i] == '%' {
			if i+1 < len(fm) && fm[i+1] == '%' {
				i++ // skip it
				continue
			} else {
				expected++
			}
		}
	}
	if n != expected {
		return fmt.Errorf("short read: whole %d, read %d (fmt %q scanning %q)",
			len(s), n, fm, s)
	}
	return nil
}

func readMap(lpc *lpc.LineParserContext) (*Map, error) {
	log.Printf("readMap")
	var m Map

	if err := asscanf(lpc.Current(), "%s map:", &m.Desc); err != nil {
		return nil, lpc.Wrap("banner scanf", err)
	}
	lpc.Next()

	for !lpc.EOF() {
		var ml MapLine
		line := lpc.Current()
		if line == "" {
			break
		}
		err := asscanf(line, "%d %d %d", &ml.OutputStart, &ml.InputStart,
			&ml.Length)
		if err != nil {
			return nil, lpc.Wrap("looping asscanf", err)
		}
		m.MapLines = append(m.MapLines, ml)
		lpc.Next()
	}

	return &m, nil
}

func loadInput() (*InputFile, error) {
	r := &InputFile{Maps: map[string]*Map{}}
	lpc := lpc.New(ick.Must(argv.ReadChompAll()))

	if lpc.Current()[0:7] != "seeds: " {
		return nil, lpc.Wrap("reading seeds", errors.New("bad seeds prefix"))
	}

	r.Seeds = ick.MapSlice(func(s string) int {
		return ick.Must(strconv.Atoi(s))
	}, strings.Split(lpc.Current()[7:], " "))

	lpc.Next()
	lpc.EatBlankLine()

	readMapCalled := func(i *InputFile, f, t string) error {
		log.Printf("read map from %s to %s", f, t)
		m, err := readMap(lpc)
		if err != nil {
			return err
		}
		expected := fmt.Sprintf("%s-to-%s", f, t)
		if expected != m.Desc {
			return fmt.Errorf("expected map %q doesn't match read map %q",
				expected, m.Desc)
		}
		lpc.EatBlankLine()
		m.From = f
		m.To = t
		i.Maps[m.Desc] = m
		log.Printf("read map: %+v", m)
		return nil
	}

	if err := readMapCalled(r, "seed", "soil"); err != nil {
		return nil, err
	}
	if err := readMapCalled(r, "soil", "fertilizer"); err != nil {
		return nil, err
	}
	if err := readMapCalled(r, "fertilizer", "water"); err != nil {
		return nil, err
	}
	if err := readMapCalled(r, "water", "light"); err != nil {
		return nil, err
	}
	if err := readMapCalled(r, "light", "temperature"); err != nil {
		return nil, err
	}
	if err := readMapCalled(r, "temperature", "humidity"); err != nil {
		return nil, err
	}
	if err := readMapCalled(r, "humidity", "location"); err != nil {
		return nil, err
	}

	log.Printf("input file %+v", r)

	return r, nil
}

func searchMap(target int, m *Map) int {
	for _, line := range m.MapLines {
		if target >= line.InputStart && target < line.InputStart+line.Length {
			// Found it.
			offset := target - line.InputStart
			r := offset + line.OutputStart
			// log.Printf("found %d in map %s (%+v) returning %d", target, m.Desc, line, r)
			return r
		}
	}

	// log.Printf("can't find it in map %s, returning input", m.Desc)
	return target
}

func search(seed int, in *InputFile) int {
	order := []string{"seed", "soil", "fertilizer", "water",
		"light", "temperature", "humidity", "location", "SENTRY",
	}

	target := seed
	for i := 0; order[i] != "location"; i++ {
		m := in.Maps[fmt.Sprintf("%s-to-%s", order[i], order[i+1])]
		target = searchMap(target, m)
		// log.Printf("new target %d", target)
	}

	return target
}

func expandSeeds(in *InputFile) []int {
	r := []int{}
	for i := 0; i < len(in.Seeds); i += 2 {
		start := in.Seeds[i]
		end := in.Seeds[i] + in.Seeds[i+1]
		for j := start; j < end; j++ {
			r = append(r, j)
		}
	}
	log.Printf("%d expanded seeds", len(r))
	return r
}

type seedAndID struct {
	seed int
	id   int
}

func searchSeeds(in *InputFile) int {
	expandedSeeds := expandSeeds(in)
	log.Printf("%d expandedSeeds", len(expandedSeeds))

	if doShuffle {
		log.Printf("shuffle...")
		rand.Shuffle(len(expandedSeeds), func(i, j int) {
			expandedSeeds[i], expandedSeeds[j] = expandedSeeds[j], expandedSeeds[i]
		})
	}

	ch := make(chan seedAndID, 1024)
	go func() {
		log.Printf("writing work items to channel...")
		for id, seed := range expandedSeeds {
			ch <- seedAndID{seed: seed, id: id}
		}
		close(ch)
	}()

	bests := make(chan int, 128)

	log.Printf("start children...")
	wg := sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(ii int) {
			best := 9999999999999
			for item := range ch {
				loc := search(item.seed, in)
				if loc < best {
					log.Printf("child %d found better seed %d (%d) at location %d",
						ii, item.seed, item.id, loc)
					best = loc
					bests <- best
				}
			}
			wg.Done()
		}(i)
	}

	finalBest := make(chan int, 1)
	go func() {
		log.Printf("collecting bests")
		best := <-bests
		log.Printf("initial best %d", best)
		for better := range bests {
			if better < best {
				log.Printf("better!: %d", better)
				best = better
			}
		}
		log.Printf("final best %d", best)
		finalBest <- best
	}()

	log.Printf("awaiting spawn")
	wg.Wait()
	close(bests) // no more bests

	return <-finalBest
}

func main() {
	in, err := loadInput()
	if err != nil {
		log.Fatalf("can't load: %v", err)
	}
	location := searchSeeds(in)
	log.Printf("min location %d", location)
}
