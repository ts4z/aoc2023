package main

/*
 * Data-sharing issues are a problem. Resetting state is tricky here,
 * since there are two data structures to keep in sync.  For part1, it's
 * immutable, so no problem.  For part 2 I had bugs for a while.
 * I think they're fixed, but I'm not sure.
 *
 * This uses fixed-length arrays out of convenience.
 */

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

type BrickSpace = [10][10][300]int

func CopySpace(dst BrickSpace, src BrickSpace) {
	for i := 0; i < len(space); i++ {
		for j := 0; j < len(space[0]); j++ {
			for k := 0; k < len(space[0][0]); k++ {
				dst[i][j][k] = src[i][j][k]
			}
		}
	}
}

func DeepCloneBrickMap(bm map[int]*Brick) map[int]*Brick {
	noo := map[int]*Brick{}
	for k, v := range bm {
		nv := &Brick{}
		*nv = *v
		noo[k] = nv
	}
	return noo
}

var (
	space    BrickSpace
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
				// log.Printf("space[%d][%d][%d] = %d", i, j, k, b.Id)
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

/*
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
*/

var chainFallersCache = map[int][]int{}

func (b *Brick) ChainFallers() []int {
	if cached, ok := chainFallersCache[b.Id]; ok {
		return cached
	}

	answers := map[int]bool{}
	tocheck := make(chan int, 1000)
	checked := map[int]bool{}

	tocheck <- b.Id

	for len(tocheck) > 0 {
		checking := <-tocheck

		if _, ok := checked[checking]; ok {
			// did it
			continue
		} else {
			checked[checking] = true
		}

		more := brickMap[checking].SolelySupports()
		for _, faller := range more {
			tocheck <- faller
			answers[faller] = true
		}
	}

	// don't chain delete myself
	delete(answers, b.Id)

	ans := ick.Keys(answers)
	chainFallersCache[b.Id] = ans
	return ans
}

func (b *Brick) SolelySupports() []int {
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

	fallers := []int{}

	// log.Printf("brick %d is supporting %d bricks", b.Id, len(overs))
	for over := range overs {
		if len(brickMap[over].BricksBelow()) == 1 {
			fallers = append(fallers, over)
		}
	}
	return fallers
}

func dropBricks() int {
	drops := map[int]bool{}
	dropsThisPass := 1
	for dropsThisPass > 0 {
		dropsThisPass = 0
		bricks := ick.Values(brickMap)
		sort.Slice(bricks, func(i, j int) bool {
			return LessByZ(bricks[i], bricks[j])
		})

		for _, b := range bricks {
			for b.From.Z > 1 && b.IsEmptyBelow() {
				drops[b.Id] = true
				dropsThisPass++
				b.Drop1()
			}
		}

		log.Printf("%d drops this pass", dropsThisPass)
	}
	return len(drops)
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
	dropBricks()

	log.Printf("counting bricks")
	unsupporting := 0
	for _, brick := range brickMap {
		if len(brick.SolelySupports()) == 0 {
			log.Printf("%v IsDisintegratable", brick.Id)
			unsupporting++
		}
	}
	fmt.Printf("part1: %d\n", unsupporting)

	// BF&I
	drops := 0
	var backupSpace BrickSpace = space
	backupBrickMap := DeepCloneBrickMap(brickMap)
	log.Printf("backupBrickMap %p", backupBrickMap)
	log.Printf("brickMap %p", brickMap)
	for id, b := range backupBrickMap {
		space = backupSpace
		brickMap = DeepCloneBrickMap(backupBrickMap)
		b.RemoveBrickInSpace()
		n := dropBricks()
		log.Printf("removing brick %d dropped %d", id, n)
		drops += n
	}

	fmt.Printf("part2: %d\n", drops)
}
