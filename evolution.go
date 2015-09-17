//+build ignore

// have boards evolve and extract statistics
package main

import (
	"fmt"
	"math/rand"

	. "."
)

func main() {
	N := 512
	b := MakeBoard(N, N)

	gens := 1200
	res := 1
	for f := 1; f < 1000; f += 5 {
		SetRand(b, int64(f), rand.Float64())
		for i := 0; i < gens; i += res {
			fmt.Println(i, avg(b))
			b.Advance(res)
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
