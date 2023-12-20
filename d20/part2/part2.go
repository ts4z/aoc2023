package main

// Another LCM puzzle, another cheat day for me, with a minor spoiler from
// Reddit.
//
// part1.go actually seems to implement this correctly but part2.go also will
// print out the cycle lengths that can be LCMed (but of course they're prime).

import (
	"fmt"
	"log"
	"strings"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

var dispatcher *Dispatcher
var button *Button = &Button{}
var rx *Rx = &Rx{}

type PulseType int

func (pt PulseType) String() string {
	return map[PulseType]string{
		Low:  "Low",
		High: "High",
	}[pt]
}

const (
	Low  = 1
	High = 2
)

type Message struct {
	Pulse PulseType
	From  string
	To    string
}

type Module interface {
	Name() string
	Receive(*Message)
}

type FlipFlop struct {
	name string
	on   bool
}

func (ff *FlipFlop) Name() string {
	return ff.name
}

func (ff *FlipFlop) Receive(m *Message) {
	switch m.Pulse {
	case High:
		return
	default:
		log.Fatalf("unknown pulse type %d", m.Pulse)
	case Low:
		if ff.on {
			ff.on = false
			dispatcher.SendPulse(ff, Low)
		} else {
			ff.on = true
			dispatcher.SendPulse(ff, High)
		}
	}
}

func NewFlipFlop(dests []string) *FlipFlop {
	return &FlipFlop{
		on: false,
	}
}

type Broadcaster struct{}

func (b *Broadcaster) Name() string {
	return "broadcaster"
}

func (b *Broadcaster) Receive(m *Message) {
	dispatcher.SendPulse(b, m.Pulse)
}

type Conjunction struct {
	name   string
	states map[string]PulseType
}

func (c *Conjunction) Name() string {
	return c.name
}

type Output struct{}

func (o *Output) Name() string {
	return "output"
}

func (o *Output) Receive(m *Message) {
	// log.Printf("output module received message %+v", m)
}

var (
	cyclePeriod   = map[string]int{}
	cycleReported = map[string]bool{}
)

func (c *Conjunction) Receive(m *Message) {
	c.states[m.From] = m.Pulse
	// log.Printf("state of %q: %+v", c.name, c.states)
	lows := ick.Grep(func(pt PulseType) bool { return pt == Low }, ick.Values(c.states))
	if len(lows) == 0 {
		// all high, send low
		dispatcher.SendPulse(c, Low)
	} else {
		n := c.Name()
		if n == "qn" || n == "xf" || n == "xn" || n == "zl" {
			presses := button.Presses()
			// log.Printf("&%s high at %d", n, presses)
			if last, ok := cyclePeriod[n]; !ok {
				cyclePeriod[n] = presses
			} else if !cycleReported[n] {
				cycleReported[n] = true
				log.Printf("%s cycle %d", n, presses-last)
			}
		}
		// at least one low, send high
		dispatcher.SendPulse(c, High)
	}
}

type Dispatcher struct {
	queue      chan *Message
	modules    map[string]Module
	srcToDst   map[string]map[string]struct{}
	dstToSrc   map[string]map[string]struct{}
	pulsesSent map[PulseType]int
}

func (d *Dispatcher) SendPulse(fromModule Module, pt PulseType) {
	src := fromModule.Name()
	for _, dst := range ick.Keys(d.srcToDst[src]) {
		d.pulsesSent[pt]++
		// log.Printf("SEND: %s -%v-> %s", fromModule.Name(), pt, dst)
		m := &Message{
			From:  src,
			To:    dst,
			Pulse: pt,
		}

		select {
		case d.queue <- m:
		default:
			log.Fatalf("dispatcher queue full")
		}
	}
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		queue:      make(chan *Message, 100000),
		srcToDst:   map[string]map[string]struct{}{},
		dstToSrc:   map[string]map[string]struct{}{},
		modules:    map[string]Module{},
		pulsesSent: map[PulseType]int{},
	}
}

func (d *Dispatcher) Wire(m Module, destinations []string) {
	d.modules[m.Name()] = m
	for _, dst := range destinations {
		if _, ok := d.srcToDst[m.Name()]; !ok {
			d.srcToDst[m.Name()] = map[string]struct{}{}
		}
		d.srcToDst[m.Name()][dst] = struct{}{}
		if _, ok := d.dstToSrc[dst]; !ok {
			d.dstToSrc[dst] = map[string]struct{}{}
		}
		d.dstToSrc[dst][m.Name()] = struct{}{}
	}
}

func (d *Dispatcher) DispatchUntilEmpty() {
	for {
		if len(d.queue) == 0 {
			return
		}

		m := <-d.queue
		// log.Printf("dispatch %+v", m)
		if dm, ok := d.modules[m.To]; ok {
			dm.Receive(m)
		} else {
			log.Printf("message to %q dropped: %+v", m.To, m)
		}
	}
}

func (d *Dispatcher) SourcesOf(name string) []string {
	return ick.Keys(d.dstToSrc[name])
}

type Button struct {
	presses int
}

func (b *Button) Presses() int {
	return b.presses
}

func (b *Button) Press() {
	b.presses++
	dispatcher.SendPulse(b, Low)
}

func (b *Button) Name() string {
	return "button"
}

func (b *Button) Receive(m *Message) {
	log.Fatalf("buttons do not receive messages like %+v", m)
}

type Rx struct {
	lows []int
}

func (rx *Rx) Name() string {
	return "rx"
}

func (rx *Rx) Receive(m *Message) {
	// log.Printf("rx received %v at press %d", m.Pulse, button.Presses())
	if m.Pulse == Low {
		rx.lows = append(rx.lows, button.Presses())
	}
}

func main() {
	log.Printf("hello")
	lines := ick.Must(argv.ReadChompAll())
	log.Printf("snarfed lines")

	dispatcher = NewDispatcher()

	// used in an example
	dispatcher.Wire(&Output{}, []string{})
	dispatcher.Wire(button, []string{"broadcaster"})
	dispatcher.Wire(&Rx{}, []string{})

	var conjunctions []*Conjunction

	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		moduleName, right := parts[0], parts[1]

		destinations := strings.Split(right, ", ")

		var newModule Module

		if moduleName == "broadcaster" {
			newModule = &Broadcaster{}
		} else if moduleName[0] == '%' {
			moduleName = moduleName[1:]
			newModule = &FlipFlop{
				name: moduleName,
				on:   false,
			}
		} else if moduleName[0] == '&' {
			moduleName := moduleName[1:]
			c := &Conjunction{
				name:   moduleName,
				states: map[string]PulseType{},
			}
			newModule = c
			conjunctions = append(conjunctions, c)
		} else {
			log.Fatalf("can't parse left side %q in line %q", moduleName, line)
		}

		dispatcher.Wire(newModule, destinations)
	}

	for _, cm := range conjunctions {
		srcs := dispatcher.SourcesOf(cm.Name())
		for _, src := range srcs {
			cm.states[src] = Low
		}
	}

	for k, v := range dispatcher.modules {
		log.Printf("module: %s %+v", k, v)
	}

	for i := 0; i < 1000; i++ {
		// log.Printf("---- pushing button ----")
		button.Press()
		dispatcher.DispatchUntilEmpty()
	}

	fmt.Printf("%d low pulses\n", dispatcher.pulsesSent[Low])
	fmt.Printf("%d high pulses\n", dispatcher.pulsesSent[High])
	fmt.Printf("%d\n", ick.Product(ick.Values(dispatcher.pulsesSent)))

	for len(rx.lows) == 0 {
		// if button.Presses()&0xFFF == 0xFFF {
		// 	log.Printf("presses: %d", button.Presses())
		// }
		button.Press()
		dispatcher.DispatchUntilEmpty()
	}

	fmt.Printf("first rx low at %d (%d)", rx.lows[0], len(rx.lows))
}
