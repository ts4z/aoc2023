package lpc

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

type LineParserContext struct {
	lines       []string
	currentLine int
}

func Argv() *LineParserContext {
	return New(ick.Must(argv.ReadChompAll()))
}

// New gets ya a new one.  Load data with argv package and shove it in here.
// Now you have a primitive way to snarf the file and print line numbers.
func New(lines []string) *LineParserContext {
	return &LineParserContext{
		lines:       lines,
		currentLine: 0,
	}
}

func (c *LineParserContext) LineNumber() int {
	// lines count from 1 for humans, 0 for computers
	return c.currentLine + 1
}

func (c *LineParserContext) EOF() bool {
	return c.currentLine >= len(c.lines)
}

func (c *LineParserContext) Next() {
	c.currentLine++
}

func (c *LineParserContext) Current() string {
	if c.currentLine <= len(c.lines) {
		return c.lines[c.currentLine]
	} else {
		return ""
	}
}

func (c *LineParserContext) Wrap(when string, err error) error {
	return fmt.Errorf("line %d (%q) %s: %w", c.LineNumber(), c.Current(),
		when, err)
}

func (c *LineParserContext) MustEatBlankLine() {
	if err := c.EatBlankLine(); err != nil {
		log.Fatalf("can't parse blank line: %v", err)
	}
}

func (c *LineParserContext) EatBlankLine() error {
	if c.EOF() {
		// Tolerate this, as some AOC inputs lack a terminating blank line
		// on the last block.
		return nil
	}
	if c.Current() == "" {
		c.Next()
		return nil
	}
	return fmt.Errorf("line %d: expected blank line, saw %q", c.LineNumber(),
		c.Current())
}
