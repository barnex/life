package life

import "fmt"

type Board struct {
	rows, cols int
	quads      [][]uint64
}

func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

func (b *Board) advance() {
}

func nextState(alive byte, neighbors byte) byte {
	if alive == 1 && neighbors == 2 || neighbors == 3 {
		return 1
	}
	return 0
}

func (b *Board) Set(R, C int, v bool) {
	var val uint64
	if v {
		val = 1
	}

	r, w, s := rowWordShift(R, C)
	b.quads[r][w] = setQuad(b.quads[r][w], s, val)
}

func (b *Board) Get(R, C int) bool {
	if R < 0 || C < 0 || R >= b.Rows() || C >= b.Cols() {
		return false
	}

	r, w, s := rowWordShift(R, C)
	word := getQuad(b.quads[r][w], s)

	switch word {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(fmt.Sprintf("bad quad at %v,%v: %v", R, C, word))
	}
}

func rowWordShift(row, col int) (r uint, w int, s uint) {
	r, w, s = uint(row), col/16, (uint(col)%16)*4
	return
}

func setQuad(word uint64, s uint, val uint64) uint64 {
	word &= (0xFFFFFFFFFFFFFFFF ^ (0xE << s))
	word |= (val << s)
	return word
}

func getQuad(word uint64, s uint) uint64 {
	word >>= s
	word &= 0xF
	return word
}

func (b *Board) Rows() int {
	return b.rows
}

func (b *Board) Cols() int {
	return b.cols
}

func MakeBoard(rows, cols int) *Board {

	ints := ((cols + 15) / 16)

	return &Board{
		rows:  rows,
		cols:  cols,
		quads: makeMatrix(rows, ints),
	}
}

func makeMatrix(rows, cols int) [][]uint64 {
	all := make([]uint64, rows*cols)
	c := make([][]uint64, rows)
	for i := range c {
		c[i] = all[i*cols : (i+1)*cols]
	}
	return c
}
