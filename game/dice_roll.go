package game

import (
	"bytes"
	"fmt"
	"math/rand"
)

const (
	DOT rune = '\u25CF' // ●
)

// Valid Values are from 1-6
type Die uint8
type Dice [2]Die

/**
 * Does a pseudo dice roll of 2 6 sided dice
 *
 * @param r Random number generator for rolling the dice
 *
 * @return A copy of the Dice Roll
 */
func (d *Dice) roll(r *rand.Rand) Dice {
	d[0] = Die(r.Intn(5) + 1)
	d[1] = Die(r.Intn(5) + 1)
	return *d
}

/**
 * Converts a set of rolled dice to points in the game
 *
 * Rolls 1-3 are as follows:
 * - 2-6 & 8-12 Points as is
 * - 7 is 70 Points
 *
 * Rolls 4+ are as follows:
 * - Same number twice doubles the current points
 * - 2-6 & 8-12 not doubles Points as is
 * - 7 ends the round
 *
 * @param rollNum Roll number in round
 * @param currPoints Running total of points thus far
 *
 * @return Running Point Total and true to keep rolling for the round, otherwise false
 */
func (d Dice) Points(rollNum int, currPoints uint) (uint, bool) {
	num := d[0] + d[1]
	points := uint(num)

	// Safe Rolls 1-3
	if rollNum <= 3 {
		if num == 7 {
			points = 70
		}
		return points + currPoints, true
	}

	// Roll 4+ and Done
	if num == 7 {
		return 0, false
	}

	// Doubles !!!
	if d[0] == d[1] {
		return currPoints * 2, true
	}

	// Just another roll
	return points + currPoints, true
}

func (d Dice) String() string {
	//     d  r  c
	var c [2][3][3]rune

	// Set all indices to spaces
	for die := 0; die < 2; die++ {
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				c[die][row][col] = ' '
			}
		}
	}

	// Change spaces to does depending on the number
	for idx, num := range d {
		switch num {
		case 1:
			c[idx][1][1] = DOT
		case 2:
			c[idx][0][0] = DOT
			c[idx][2][2] = DOT
		case 3:
			c[idx][0][0] = DOT
			c[idx][1][1] = DOT
			c[idx][2][2] = DOT
		case 4:
			c[idx][0][0] = DOT
			c[idx][0][2] = DOT
			c[idx][2][0] = DOT
			c[idx][2][2] = DOT
		case 5:
			c[idx][0][0] = DOT
			c[idx][0][2] = DOT
			c[idx][1][1] = DOT
			c[idx][2][0] = DOT
			c[idx][2][2] = DOT
		case 6:
			c[idx][0][0] = DOT
			c[idx][0][2] = DOT
			c[idx][1][0] = DOT
			c[idx][1][2] = DOT
			c[idx][2][0] = DOT
			c[idx][2][2] = DOT
		}
	}

	// Format the dice
	buf := bytes.NewBufferString("")
	fmt.Fprintf(buf, "┌─────────┐ ┌─────────┐\n")
	fmt.Fprintf(buf, "│ %c     %c │ │ %c     %c │\n", c[0][0][0], c[0][0][2], c[1][0][0], c[1][0][2])
	fmt.Fprintf(buf, "│ %c  %c  %c │ │ %c  %c  %c │\n", c[0][1][0], c[0][1][1], c[0][1][2], c[1][1][0], c[1][1][1], c[1][1][2])
	fmt.Fprintf(buf, "│ %c     %c │ │ %c     %c │\n", c[0][2][0], c[0][2][2], c[1][2][0], c[1][2][2])
	fmt.Fprintf(buf, "└─────────┘ └─────────┘\n")

	return buf.String()
}
