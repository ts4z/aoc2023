package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/ts4z/aoc2023/argv"
	"github.com/ts4z/aoc2023/ick"
)

var brickRE = regexp.MustCompile(`^(\d+),(\d+),(\d+)~(\d+),(\d+),(\d+)$`)

type Cube struct{ X, Y, Z int }
type Brick struct {
	Id       int
	From, To Cube
}

var (
	space    = [10][10][300]int{}
	brickMap = map[int]*Brick{}
)

func printSpace() {
	for k := len(space[0][0]) - 1; k > 0; k-- {
		fmt.Printf(" %3d ----------------------------\n", k)
		for i := 0; i < len(space); i++ {
			fmt.Printf("[%4d]", k)
			for j := 0; j < len(space[0]); j++ {
				fmt.Printf(".%4d", space[i][j][k])
			}
			fmt.Printf(".\n")
		}
	}
}

func cmpByFields[T any](accessors []func(T) int, left, right T) int {
	for _, accessor := range accessors {
		lv := accessor(left)
		rv := accessor(right)
		if lv == rv {
			continue // TBD
		}
		if lv < rv {
			return -1
		}
		return 1
	}
	return 0
}

func LessByZ(left, right *Brick) bool {
	return cmpByFields([]func(*Brick) int{
		func(b *Brick) int { return b.From.Z },
		func(b *Brick) int { return b.To.Z },
		func(b *Brick) int { return b.From.X },
		func(b *Brick) int { return b.From.Y },
		func(b *Brick) int { return b.To.X },
		func(b *Brick) int { return b.To.Y },
	}, left, right) < 0
}

func GreaterByZ(left, right *Brick) bool {
	return cmpByFields([]func(*Brick) int{
		func(b *Brick) int { return b.From.Z },
		func(b *Brick) int { return b.To.Z },
		func(b *Brick) int { return b.From.X },
		func(b *Brick) int { return b.From.Y },
		func(b *Brick) int { return b.To.X },
		func(b *Brick) int { return b.To.Y },
	}, left, right) > 0
}

func (b *Brick) EnsureLowZFirst() {
	if b.From.Z > b.To.Z {
		b.From, b.To = b.To, b.From
	}
}

func (b *Brick) Drop1() {
	b.RemoveBrickInSpace()
	b.From.Z--
	b.To.Z--
	b.InsertBrickInSpace()
}

func (b *Brick) String() string {
	return fmt.Sprintf("Brick{Id:%d,{%d,%d,%d},{%d,%d,%d}}",
		b.Id, b.From.X, b.From.Y, b.From.Z,
		b.To.X, b.To.Y, b.To.Z)
}

func (b *Brick) InsertBrickInSpace() {
	// log.Printf("Inserting %v", b)
	for i := b.From.X; i <= b.To.X; i++ {
		for j := b.From.Y; j <= b.To.Y; j++ {
			for k := b.From.Z; k <= b.To.Z; k++ {
				occ := space[i][j][k]
				if occ != 0 {
					log.Fatalf("can't insert brick %v in space at %d,%d,%d, occupied by %d (%v)",
						b, i, j, k, occ, brickMap[occ])
				}
				log.Printf("space[%d][%d][%d] = %d", i, j, k, b.Id)
				space[i][j][k] = b.Id
			}
		}
	}
}

func (b *Brick) RemoveBrickInSpace() {
	// log.Printf("Removing %v", b)
	for i := b.From.X; i <= b.To.X; i++ {
		for j := b.From.Y; j <= b.To.Y; j++ {
			for k := b.From.Z; k <= b.To.Z; k++ {
				occ := space[i][j][k]
				if occ != b.Id {
					log.Fatalf("can't remove brick %v in space at %d,%d,%d, occupied by %d (%v)",
						b, i, j, k, occ, brickMap[occ])
				}
				space[i][j][k] = 0
			}
		}
	}
}

func (b *Brick) IsEmptyBelow() bool {
	if b.From.Z == 1 {
		return false
	}

	k := b.From.Z

	for i := b.From.X; i <= b.To.X; i++ {
		for j := b.From.Y; j <= b.To.Y; j++ {
			under := space[i][j][k-1]
			if under == 0 {
				continue
			}
			if under == b.Id {
				return true // will be OK when brick moves
			}
			return false
		}
	}
	return true
}

func (b *Brick) BricksBelow() []int {
	if b.From.Z == 1 {
		log.Printf("brick %d on ground", b.Id)
		return nil
	}

	below := map[int]bool{}

	k := b.From.Z

	for i := b.From.X; i <= b.To.X; i++ {
		for j := b.From.Y; j <= b.To.Y; j++ {
			under := space[i][j][k-1]
			if under != 0 {
				below[under] = true
			}
		}
	}
	log.Printf("bricksBelow %d are %+v", b.Id, below)
	return ick.Keys(below)
}

func (b *Brick) IsDisintegratable() bool {
	overs := map[int]bool{}
	k := b.To.Z
	for i := b.From.X; i <= b.To.X; i++ {
		for j := b.From.Y; j <= b.To.Y; j++ {
			over := space[i][j][k+1]
			if over != 0 {
				overs[over] = true
			}
		}
	}

	log.Printf("brick %d is supporting %d bricks", b.Id, len(overs))
	for over := range overs {
		if len(brickMap[over].BricksBelow()) == 1 {
			// assert(BricksBelow[0] == b.Id
			// the brick we are supporting is supported only by us
			return false
		}
	}
	return true
}

func main() {
	lines := ick.Must(argv.ReadChompAll())

	log.Printf("parsing bricks")
	for i, line := range lines {
		brickId := i + 1

		ms := brickRE.FindSubmatch([]byte(line))
		if ms == nil {
			log.Fatalf("re mismatch")
		}
		m := ick.MapSlice(func(b []byte) int { return ick.Atoi(string(b)) }, ms[1:])
		brick := &Brick{brickId, Cube{m[0], m[1], m[2]}, Cube{m[3], m[4], m[5]}}
		brick.EnsureLowZFirst()
		brickMap[brickId] = brick
	}

	log.Printf("inserting bricks")
	for _, brick := range brickMap {
		brick.InsertBrickInSpace()
	}

	log.Printf("dropping bricks")
	drops := 1
	// printSpace()
	for drops > 0 {
		drops = 0
		bricks := ick.Values(brickMap)
		sort.Slice(bricks, func(i, j int) bool {
			return LessByZ(bricks[i], bricks[j])
		})

		for _, b := range bricks {
			for b.From.Z > 1 && b.IsEmptyBelow() {
				drops++
				b.Drop1()
			}
		}

		log.Printf("%d drops this pass", drops)
		// printSpace()
	}

	log.Printf("counting bricks")
	unsupporting := 0
	for _, brick := range brickMap {
		if brick.IsDisintegratable() {
			log.Printf("%v IsDisintegratable", brick.Id)
			unsupporting++
		}
	}
	fmt.Printf("%d\n", unsupporting)
}
