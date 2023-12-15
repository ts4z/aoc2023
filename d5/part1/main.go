package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
	"github.com/ts4z/aoc2023/lpc"
)

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

	r.Seeds = ick.MapSlice(ick.Atoi, strings.Split(lpc.Current()[7:], " "))

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
			log.Printf("found %d in map %s (%+v) returning %d", target,
				m.Desc, line, r)
			return r
		}
	}

	log.Printf("can't find it in map %s, returning input", m.Desc)
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
		log.Printf("new target %d", target)
	}

	return target
}

func searchSeeds(in *InputFile) map[int]int {
	r := map[int]int{}
	for _, seed := range in.Seeds {
		r[seed] = search(seed, in)
	}
	return r
}

func main() {
	in, err := loadInput()
	if err != nil {
		log.Fatalf("can't load: %v", err)
	}
	locations := searchSeeds(in)
	log.Printf("locations: %+v\n", locations)
	minimum := ick.Min(ick.Values(locations))
	log.Printf("min location %d", minimum)
}
