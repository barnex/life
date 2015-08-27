package main

import (
	"fmt"
	"time"
)

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
	b := ParseBoard(
		`......
.xxx..
......
`)

	for {
		fmt.Println(Fmt(b))
		b.Advance(1)
		time.Sleep(500 * time.Millisecond)
	}

}
