package life

import "fmt"

type Board struct {
	rows, cols int
	nibs       [][]uint64
	hsum       [][]uint64
	vsum       [][]uint64
}

func (b *Board) count() {
	for r, row := range b.nibs {
		rowHSum(b.hsum[r], row)
	}
}

func rowHSum(dst, src []uint64) {
	var prev uint64 = 0

	for i := 0; i < len(src)-1; i++ {
		centre := src[i]
		left := (prev & 0xF000000000000000) >> ((16-1)*4) | (centre << 4)
		right := src[i+1]

		dst[i] = 
	}
}

func countRow(row []uint64) []uint64 {
	sum := make([]uint64, len(row))
	for w, word := range row {
		sum[w] = word + (word << 4) + (word >> 4)
	}
	return sum
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
	b.nibs[r][w] = setnib(b.nibs[r][w], s, val)
}

func (b *Board) Get(R, C int) bool {
	if R < 0 || C < 0 || R >= b.Rows() || C >= b.Cols() {
		return false
	}

	r, w, s := rowWordShift(R, C)
	word := getnib(b.nibs[r][w], s)

	switch word {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(fmt.Sprintf("bad nib at %v,%v: %v", R, C, word))
	}
}

func rowWordShift(row, col int) (r uint, w int, s uint) {
	r, w, s = uint(row), col/16, (uint(col)%16)*4
	return
}

func setnib(word uint64, s uint, val uint64) uint64 {
	word &= (0xFFFFFFFFFFFFFFFF ^ (0xE << s))
	word |= (val << s)
	return word
}

func getnib(word uint64, s uint) uint64 {
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
		rows: rows,
		cols: cols,
		nibs: makeMatrix(rows, ints),
		hsum: makeMatrix(rows, ints),
		vsum: makeMatrix(rows, ints),
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
