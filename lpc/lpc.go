package lpc

import (
	"fmt"
)

type LineParserContext struct {
	lines       []string
	currentLine int
}

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

func (c *LineParserContext) EatBlankLine() error {
	if c.EOF() {
		// Tolerate this, as some AOC inputs lack a terminating blank line
		// on the last block.
		return nil
	}
	if c.Current() == "" {
		c.currentLine++
		return nil
	}
	return fmt.Errorf("line %d: expected blank line, saw %q", c.LineNumber(),
		c.Current())
}
