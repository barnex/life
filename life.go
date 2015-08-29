package main

import (
	"bytes"
	"fmt"
	"time"
)

func Fmt(b *Board) string {
	rows := b.Rows()
	cols := b.Cols()
	var buf bytes.Buffer
	fmt.Fprintln(&buf)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := b.Get(i, j)
			if v {
				fmt.Fprint(&buf, "x")
			} else {
				fmt.Fprint(&buf, " ")
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

func ParseBoard(in [][]bool) *Board {
	rows := len(in)
	cols := len(in[0])
	b := MakeBoard(rows, cols)

	for r, row := range in {
		for c, v := range row {
			b.Set(r, c, v)
		}
	}
	return b
}

func NextState(alive bool, neighbors int) bool {
	if alive {
		return neighbors == 2 || neighbors == 3
	} else {
		return neighbors == 3
	}
}

func makeMatrix(rows, cols int) [][]bool {
	all := make([]bool, rows*cols)
	c := make([][]bool, rows)
	for i := range c {
		c[i] = all[i*cols : (i+1)*cols]
	}
	return c
}

func (b *Board) wrap(r, c int) (int, int) {
	r = (r + b.Rows()) % b.Rows()
	c = (c + b.Cols()) % b.Cols()
	return r, c
}

func main() {
	const (
		X = true
		O = false
	)
	b := ParseBoard(
		[][]bool{
			{O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{X, X, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, X, X, X, X, X, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		})

	for {
		fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n")
		fmt.Println(Fmt(b))
		b.Advance(1)
		time.Sleep(500 * time.Millisecond)
	}

}
