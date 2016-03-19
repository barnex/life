/*
Command evolution runs game of life from random initial states with varying
fill fractions (0-100%), and outputs how the state evolves towards equilibrium.
Output is two-column table:
	generation, fill fraction
*/
package main

import (
	"fmt"

	. "github.com/barnex/life"
)

func main() {
	N := 512
	b := MakeBoard(N, N)

	gens := 10000
	for f := 1.; f < 100; f += 1 {
		SetRand(b, 0, f/100)
		res := 1
		for i := 0; i < gens; i += res {
			fmt.Println(i, avg(b))
			b.Advance(res)
			if i > 100 {
				res = i / 100
			}
		}
		fmt.Println()
	}
}

func avg(b *Board) float64 {
	rows := b.Rows()
	cols := b.Cols()

	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if b.Get(r, c) {
				count++
			}
		}
	}

	return float64(count) / float64(rows*cols)
}
