package main

import (
	"fmt"
	"log"

	"github.com/ts4z/aoc2023/aoc"
)

type Position struct {
	Row    int
	Column int
}

func (pos Position) String() string {
	return fmt.Sprintf("%#v", pos)
}

type Direction struct {
	Name string
	Step func(Position) Position
}

func GoNorth(p Position) Position {
	return Position{p.Row - 1, p.Column}
}
func GoWest(p Position) Position {
	return Position{p.Row, p.Column - 1}
}
func GoEast(p Position) Position {
	return Position{p.Row, p.Column + 1}
}
func GoSouth(p Position) Position {
	return Position{p.Row + 1, p.Column}
}

var North = &Direction{"North", GoNorth}
var South = &Direction{"South", GoSouth}
var East = &Direction{"East", GoEast}
var West = &Direction{"West", GoWest}

type Ray struct {
	Position  Position
	Direction *Direction
}

func (ray Ray) String() string {
	return fmt.Sprintf("Ray{%v,%s}", ray.Position, ray.Direction.Name)
}

type SpaceAction func(Ray) []Ray

func (ray Ray) Bend(d *Direction) []Ray {
	newRay := []Ray{{d.Step(ray.Position), d}}
	log.Printf("ray %+v bent, now %+v", ray, newRay)
	return newRay
}

func slashMirror(ray Ray) []Ray {
	log.Printf("slashMirror %v", ray)
	switch ray.Direction {
	case North:
		return ray.Bend(East)
	case West:
		return ray.Bend(South)
	case East:
		return ray.Bend(North)
	case South:
		return ray.Bend(West)
	default:
		log.Fatalf("slashMirror can't handle ray %+v", ray)
		return nil
	}
}

func backslashMirror(ray Ray) []Ray {
	log.Printf("backslashMirror %v", ray)
	switch ray.Direction {
	case North:
		return ray.Bend(West)
	case West:
		return ray.Bend(North)
	case East:
		return ray.Bend(South)
	case South:
		return ray.Bend(East)
	default:
		log.Fatalf("backslashMirror can't handle ray %+v", ray)
		return nil
	}
}

func dashSplitter(ray Ray) []Ray {
	switch ray.Direction {
	case East:
		return passthru(ray)
	case West:
		return passthru(ray)
	case North:
		fallthrough
	case South:
		return []Ray{
			{ray.Position, East},
			{ray.Position, West},
		}
	default:
		log.Fatalf("dashSplitter can't handle ray %+v", ray)
		return nil
	}
}

func pipeSplitter(ray Ray) []Ray {
	switch ray.Direction {
	case East:
		fallthrough
	case West:
		return []Ray{
			{ray.Position, North},
			{ray.Position, South},
		}
	case North:
		return passthru(ray)
	case South:
		return passthru(ray)
	default:
		log.Fatalf("dashSplitter can't handle ray %+v", ray)
		return nil
	}
}

func CheckBounds(a [][]byte, pos Position) bool {
	if pos.Row < 0 {
		return false
	}
	if pos.Row >= len(a) {
		return false
	}
	if pos.Column < 0 {
		return false
	}
	if pos.Column >= len(a[0]) {
		return false
	}
	return true // ok
}

func passthru(ray Ray) []Ray {
	return []Ray{{ray.Direction.Step(ray.Position), ray.Direction}}
}

var byteToAction = map[byte]SpaceAction{
	'/':  slashMirror,
	'\\': backslashMirror,
	'-':  dashSplitter,
	'|':  pipeSplitter,
	'.':  passthru,
}

func main() {
	a := aoc.ReadInputAsByteMatrix()
	fmt.Printf("map:\n")
	aoc.PrintByteMatrix(a)

	todo := []Ray{{Position{0, 0}, East}}
	done := map[string]struct{}{}
	active := map[string]struct{}{}

	for len(todo) > 0 {
		ray := todo[len(todo)-1]
		todo = todo[:len(todo)-1]

		if !CheckBounds(a, ray.Position) {
			log.Printf("ray %v out of bounds", ray)
			continue
		}

		if _, ok := done[ray.String()]; ok {
			log.Printf("already did %v", ray)
			continue
		}

		active[ray.Position.String()] = struct{}{}
		done[ray.String()] = struct{}{}

		ch := a[ray.Position.Row][ray.Position.Column]
		action := byteToAction[ch]
		newActions := action(ray)
		log.Printf("ray %v action '%c' => %v", ray, ch, newActions)
		todo = append(todo, newActions...)
	}

	fmt.Printf("%d items in active map\n", len(active))
}
