package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

func hash(s string) int {
	h := 0
	for _, ch := range []byte(s) {
		h += int(ch)
		h *= 17
		h %= 256
	}
	fmt.Printf("hash of %s = %d\n", s, h)
	return h
}

func parseInt(s string) *int {
	if len(s) == 0 {
		return nil
	} else {
		v := ick.Must(strconv.Atoi(s))
		return &v
	}
}

type boxItem struct {
	label         string
	focusingPower int
}

func printBoxes(boxes [256][]boxItem) {
	for boxNumber0, boxItems := range boxes {
		if len(boxItems) == 0 {
			continue
		}
		fmt.Printf("Box %d:", boxNumber0)
		for _, item := range boxItems {
			fmt.Printf(" [%s %d]", item.label, item.focusingPower)
		}
		fmt.Printf("\n")
	}
}

func main() {
	boxes := [256][]boxItem{}

	lines, err := argv.ReadChompAll()
	if err != nil {
		log.Fatalf("can't read: %v", err)
	}

	insts := strings.Split(lines[0], ",")

	for instructionNumber, inst := range insts {
		p := 0
		fmt.Printf("parsing %s\n", inst)
		for {
			if inst[p] == '-' || inst[p] == '=' {
				break
			}
			p++
		}
		label := inst[:p]
		boxNumber := hash(label)
		op := inst[p]
		arg := parseInt(inst[p+1:])
		sarg := ""
		if arg != nil {
			sarg = fmt.Sprintf(" %d", *arg)
		}

		fmt.Printf("i#=%d label=%s boxNumber=%d op=%c arg=%+v%s\n", instructionNumber, label, boxNumber, op, arg, sarg)

		switch op {
		case '-':
			items := boxes[boxNumber]
			boxes[boxNumber] = nil
			for _, item := range items {
				if item.label != label {
					boxes[boxNumber] = append(boxes[boxNumber], item)
				}
			}

		case '=':
			items := boxes[boxNumber]
			found := false
			for i := range items {
				if items[i].label == label {
					found = true
					items[i].focusingPower = *arg
					break
				}
			}
			if !found {
				items = append(items, boxItem{label, *arg})
			}
			boxes[boxNumber] = items
		}

		printBoxes(boxes)
	}

	answer := 0

	for boxNumber0, boxItems := range boxes {
		for position0, item := range boxItems {
			product := (boxNumber0 + 1) * (position0 + 1) * item.focusingPower
			answer += product
		}
	}

	fmt.Printf("%d\n", answer)
}
