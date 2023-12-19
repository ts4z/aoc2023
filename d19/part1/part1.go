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

type Rule struct {
	Name          string
	Conditions    []Condition
	DefaultAction string
}

type Part map[string]int

var ruleLineRE = regexp.MustCompile(`([a-z]+)\{(.*)\}`)
var ruleConditionRE = regexp.MustCompile(`([xmas])([><])(\d+):([a-z]+|[AR])`)
var partRE = regexp.MustCompile(`\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)\}`)

func parseInput(lines []string) (map[string]*Rule, []Part) {

	rules := map[string]*Rule{}
	parts := []Part{}
	lineno := 0

	for {
		line := lines[lineno]
		lineno++

		if line == "" {
			break
		}

		matches := ruleLineRE.FindSubmatch([]byte(line))
		if matches == nil {
			log.Fatalf("can't parse rule input line %q", matches)
		}
		name := string(matches[1])
		log.Printf("rule %q unparsed conditions %s", name, matches[2])
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
		log.Printf("rule %q parsed conditions %+v", name, conditions)

		rules[name] = &Rule{
			Name:          name,
			Conditions:    conditions,
			DefaultAction: unparsedConditions[len(unparsedConditions)-1],
		}

	}

	for lineno < len(lines) {
		line := lines[lineno]
		lineno++

		subs := partRE.FindSubmatch([]byte(line))
		if subs == nil {
			log.Fatalf("can't parse part %q", line)
		}
		parts = append(parts, Part{
			"x": ick.Atoi(string(subs[1])),
			"m": ick.Atoi(string(subs[2])),
			"a": ick.Atoi(string(subs[3])),
			"s": ick.Atoi(string(subs[4])),
		})
	}

	return rules, parts
}

func processPart(part Part, rules map[string]*Rule, nextRule string) string {
	rule := rules[nextRule]
	if rule == nil {
		log.Fatalf("unknown rule %q", rule)
	}

	resolveAction := func(action string) string {
		if action == "A" || action == "R" {
			log.Printf("part %v rule %s resolves to %s", rule, nextRule, action)
			return action
		} else {
			log.Printf("part %v rule %s now runs %s", rule, nextRule, action)
			return processPart(part, rules, action)
		}
	}

	for i, c := range rule.Conditions {
		rcv := c.Value
		pv := part[c.Attribute]
		match := false
		op := c.Op
		if op == '<' {
			match = pv < rcv
		} else if op == '>' {
			match = pv > rcv
		} else {
			log.Fatalf("unknown op %v in rule %+v condition %d", op, rule, i)
		}

		if match {
			return resolveAction(c.Action)
		}
	}

	return resolveAction(rule.DefaultAction)
}

func main() {
	lines := ick.Must(argv.ReadChompAll())
	rules, parts := parseInput(lines)

	for name, rule := range rules {
		log.Printf("rule %q: %+v", name, rule)
	}

	for i, part := range parts {
		log.Printf("part %d: %+v", i, part)
	}

	total := 0
	for _, part := range parts {
		action := processPart(part, rules, "in")
		if action == "A" {
			total += ick.Sum(ick.Values(part))
		} else if action != "R" {
			log.Fatalf("part %+v nonsense action %q", part, action)
		}
	}

	fmt.Printf("total %d\n", total)
}
