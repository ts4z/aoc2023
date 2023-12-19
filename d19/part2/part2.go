package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

type Op byte

type Condition struct {
	Attribute string
	Op        Op
	Value     int
	Action    string
}

func (c Condition) String() string {
	return fmt.Sprintf("Condition{%s %c %d => %s}",
		c.Attribute, c.Op, c.Value, c.Action)
}

func (c Condition) Invert() Condition {
	nc := c
	if c.Op == '<' {
		// c < n inverted is c >= n or c > n-1 if c is an int
		nc.Op = '>'
		nc.Value = c.Value - 1
	} else if c.Op == '>' {
		// c > n inverted is c <= n or c < n+1
		nc.Op = '<'
		nc.Value = c.Value + 1
	}
	nc.Action = "(inverted)"
	log.Printf("      invert %v to %v", c, nc)
	return nc
}

type Rule struct {
	Name          string
	Conditions    []Condition
	DefaultAction string
}

type Part map[string]int

var ruleLineRE = regexp.MustCompile(`([a-z]+)\{(.*)\}`)
var ruleConditionRE = regexp.MustCompile(`([xmas])([><])(\d+):([a-z]+|[AR])`)
var partRE = regexp.MustCompile(`\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)\}`)

func parseInput(lines []string) map[string]*Rule {

	rules := map[string]*Rule{}
	lineno := 0

	for {
		line := lines[lineno]
		lineno++

		if line == "" {
			return rules
		}

		matches := ruleLineRE.FindSubmatch([]byte(line))
		if matches == nil {
			log.Fatalf("can't parse rule line %d input line %q", lineno, line)
		}
		name := string(matches[1])
		// log.Printf("rule %q unparsed conditions %s", name, matches[2])
		unparsedConditions := strings.Split(string(matches[2]), ",")
		conditions := []Condition{}
		for i := 0; i < len(unparsedConditions)-1; i++ {
			cs := unparsedConditions[i]
			cms := ruleConditionRE.FindSubmatch([]byte(cs))
			if cms == nil {
				log.Fatalf("can't parse conditions substring %q of rule %q (%q)",
					cs, name, line)
			}
			conditions = append(conditions, Condition{
				Attribute: string(cms[1]),
				Op:        Op(cms[2][0]),
				Value:     ick.Atoi(string(cms[3])),
				Action:    string(cms[4]),
			})
		}
		// log.Printf("rule %q parsed conditions %+v", name, conditions)

		rules[name] = &Rule{
			Name:          name,
			Conditions:    conditions,
			DefaultAction: unparsedConditions[len(unparsedConditions)-1],
		}
	}

	return rules
}

// Cover up for some sins of data sharing.
func clone[T any](in []T) []T {
	out := make([]T, len(in))
	copy(out, in)
	return out
}

// This initially returned a result via a channel, but there's too much data
// sharing and subtle modifications in slightly bad ways; the problems are
// still here, but covered up well enough for consistent results.  (For a good time
// on a very dull Saturday night, convert this to pass its results to a channel.)
func processTree(rules map[string]*Rule) [][]Condition {
	accepts := [][]Condition{}

	var recur func(prereqs []Condition, currentRule string)

	recur = func(prereqs []Condition, currentRule string) {
		rule := rules[currentRule]

		for _, cond := range rule.Conditions {
			if cond.Action == "R" {
				// terminal dead branch
			} else if cond.Action == "A" {
				// If this rule is included, this is an accept.
				log.Printf("accept %v + %v", prereqs, cond)
				accepts = append(accepts, clone(append(prereqs, cond)))
			} else {
				// Recurse for things that match this rule.
				recur(append(prereqs, cond), cond.Action)
			}

			// Recurse for things that did not match this rule.
			inv := cond.Invert()
			prereqs = append(prereqs, inv)
		}

		def := rule.DefaultAction
		if def == "R" {
			log.Printf("default action rejects things that match: %v", prereqs)
			// done, dead branch
		} else if def == "A" {
			accepts = append(accepts, clone(prereqs))
		} else {
			recur(prereqs, def)
		}
	}

	recur([]Condition{}, "in")

	return accepts
}

func main() {
	lines := ick.Must(argv.ReadChompAll())
	rules := parseInput(lines)

	for name, rule := range rules {
		log.Printf("rule %q: %+v", name, rule)
	}

	accepts := processTree(rules)

	total := 0
	for na, reqs := range accepts {
		log.Printf("na=%d matching %+v", na, reqs)
		gt := map[string]int{"x": 0, "m": 0, "a": 0, "s": 0}
		lt := map[string]int{"x": 4001, "m": 4001, "a": 4001, "s": 4001}

		for i, c := range reqs {
			log.Printf("  i=%d %+v", i, c.String())
			if c.Op == '<' {
				lt[c.Attribute] = min(lt[c.Attribute], c.Value)
				log.Printf("    %s < %d", c.Attribute, lt[c.Attribute])
			} else if c.Op == '>' {
				gt[c.Attribute] = max(gt[c.Attribute], c.Value)
				log.Printf("    %s > %d", c.Attribute, gt[c.Attribute])
			}
		}

		cases := 1
		for _, ch := range "xmas" {
			s := string(ch)
			diff := lt[s] - gt[s] - 1
			log.Printf("%s: (%d,%d) %d values", s, gt[s], lt[s], diff)
			if diff <= 0 {
				log.Printf("dead branch %s (%d)", s, diff)
				cases = 0
			} else {
				cases *= diff
			}
		}
		total += cases
	}
	fmt.Printf("%d accept rules\n", total)
}
