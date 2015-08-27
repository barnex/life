package main

import (
	"bytes"
	"fmt"
	"strings"
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

type Board struct {
	cells [][]bool
	temp  [][]bool
}

func ParseBoard(str string) *Board {
	lines := strings.Split(str, "\n")

	rows := len(lines)
	cols := 0
	for _, l := range lines {
		if len(l) > cols {
			cols = len(l)
		}
	}

	b := MakeBoard(rows, cols)
	for r, l := range lines {
		for c := range l {
			v := l[c] == 'x'
			b.Set(r, c, v)
		}
	}
	return b
}

func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

func (b *Board) advance() {
	for i := range b.cells {
		for j := range b.cells[i] {
			neigh := b.Neighbors(i, j)
			alive := b.cells[i][j]
			b.temp[i][j] = NextState(alive, neigh)
		}
	}
	b.cells, b.temp = b.temp, b.cells
}

func (b *Board) Neighbors(r, c int) int {
	count := 0
	for rr := r - 1; rr <= r+1; rr++ {
		for cc := c - 1; cc <= c+1; cc++ {
			if b.Get(rr, cc) {
				count++
			}
		}
	}
	// do not count self
	if b.Get(r, c) {
		count--
	}
	return count
}

func (b *Board) Set(r, c int, v bool) {
	r, c = b.wrap(r, c)
	b.cells[r][c] = v
}

func (b *Board) Get(r, c int) bool {
	r, c = b.wrap(r, c)
	return b.cells[r][c]
}
func (b *Board) Rows() int {
	return len(b.cells)
}

func (b *Board) Cols() int {
	return len(b.cells[0])
}

func MakeBoard(rows, cols int) *Board {
	return &Board{
		cells: makeMatrix(rows, cols),
		temp:  makeMatrix(rows, cols),
	}
}
