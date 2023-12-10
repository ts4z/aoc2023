package esc

// package esc outputs VT102-esque escape sequences.

import (
	"fmt"
)

func Home() {
	fmt.Printf("\033[H")
}

func Clear() {
	fmt.Printf("\033[2J")
}

// Goto a row/column (note these are 1-based, that is, Home() is Goto(1,1).
func Goto(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}
