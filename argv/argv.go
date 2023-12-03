package argv

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	"github.com/ts4z/aoc2023/ick"
)

// ReadAll is the moral equivalent of using <ARGV> in a list context in perl.
func ReadAll() ([]string, error) {
	if len(os.Args) == 1 {
		return ick.ReadLines(os.Stdin)
	}

	lines := []string{}

	for i, filename := range os.Args {
		if i == 0 {
			continue
		}

		fh, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer fh.Close()

		r := bufio.NewReader(fh)

		for {
			line, err := r.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					return lines, nil
				}
				log.Printf("err %v with %d lines in buffer", err, len(lines))
				return nil, err
			}
			if line[len(line)-1] != '\n' {
				return nil, errors.New("line didn't end in LF")
			}
			lines = append(lines, line)
		}
	}

	return lines, nil
}

// ReadChompAll is the moreal equivalent of using <ARGV> in a list context,
// followed by chomping the input lines (with the canonical record separator).
func ReadChompAll() ([]string, error) {
	lines, err := ReadAll()
	if err != nil {
		return nil, err
	}
	chomped := ick.MapSlice(ick.ChompNL, lines)
	return chomped, nil
}
