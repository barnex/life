package life

import (
	"bytes"
	"fmt"
	"math/rand"
)

func BoardSet(b *Board, roff, coff int, in [][]bool) {
	for r, row := range in {
		for c, v := range row {
			b.Set(r+roff, c+coff, v)
		}
	}
}

func SetRand(b *Board, seed int64, fill float64) {
	rand.Seed(seed)
	for r := 0; r < b.Rows(); r++ {
		for c := 0; c < b.Cols(); c++ {
			v := (rand.Float64() <= fill)
			b.Set(r, c, v)
		}
	}
}

func Fmt(b *Board) string {
	rows := b.Rows()
	cols := b.Cols()
	var buf bytes.Buffer
	fmt.Fprintln(&buf)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := b.Get(i, j)
			if v {
				fmt.Fprint(&buf, "X")
			} else {
				fmt.Fprint(&buf, ".")
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}
