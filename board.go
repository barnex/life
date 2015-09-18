package life

import "unsafe"

type Board struct {
	rows, cols int
	cells      [][]byte // current cells
	temp       [][]byte // buffer for next-gen cells
	empty      []byte   // empty cell row used at borders
	colsum     []byte   //buffer for vertical sums by 3
}

func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

func (b *Board) advance() {
	rows := b.Rows()

	// first row
	r := 0
	b.advRow(r, b.empty, b.cells[r], b.cells[r+1])

	// bulk rows
	for r = 1; r < rows-1; r++ {
		b.advRow(r, b.cells[r-1], b.cells[r], b.cells[r+1])
	}

	// last rows
	r = rows - 1
	b.advRow(r, b.cells[r-1], b.cells[r], b.empty)

	// swap: temp becomes current cells
	b.cells, b.temp = b.temp, b.cells
}

func as64(bytes []byte) []uint64 {
	if len(bytes)%8 != 0 {
		panic("as64")
	}
	return (*((*[1 << 31]uint64)(unsafe.Pointer(&bytes[0]))))[:len(bytes)/8]
}

func colSum(dst, up, me, down []byte) {
	dst64 := as64(dst)
	up64 := as64(up)
	me64 := as64(me)
	down64 := as64(down)

	for i := range dst64 {
		dst64[i] = up64[i] + me64[i] + down64[i]
	}
}

func (b *Board) advRow(r int, up, me, down []byte) {

	cs := b.colsum
	colSum(cs, up, me, down)

	cols := b.Cols()
	result := b.temp[r]

	var prevCS, currCS, nextCS byte

	// first col
	c := 0
	alive := me[c]

	prevCS = 0
	currCS = cs[c]
	nextCS = cs[c+1]

	neigh := prevCS + currCS + nextCS
	result[c] = nextLUT[(alive<<4)|neigh]

	// bulk cols
	for c := 1; c < cols-1; c++ {
		alive = me[c]

		prevCS = currCS
		currCS = nextCS
		nextCS = cs[c+1]

		neigh = prevCS + currCS + nextCS
		result[c] = nextLUT[(alive<<4)|neigh]
	}

	// last col
	c = cols - 1
	alive = me[c]
	neigh = cs[c-1] + cs[c]
	result[c] = nextLUT[(alive<<4)|neigh]
}

var nextLUT [32]byte

func init() {
	for _, alive := range []byte{0, 1} {
		for neigh := byte(0); neigh <= 8; neigh++ {
			idx := (alive << 4) | neigh
			nextLUT[idx] = nextState(alive, neigh-alive) // self is included in neigh
		}
	}
}

func nextState(alive byte, neighbors byte) byte {
	if alive == 1 && neighbors == 2 || neighbors == 3 {
		return 1
	}
	return 0
}

func (b *Board) Set(r, c int, v bool) {
	if v {
		b.cells[r][c] = 1
	} else {
		b.cells[r][c] = 0
	}
}

func (b *Board) get(r, c int) byte {
	if r < 0 || c < 0 || r >= b.Rows() || c >= b.Cols() {
		return 0
	}
	return b.cells[r][c]
}
func (b *Board) Get(r, c int) bool {
	return b.get(r, c) == 1
}
func (b *Board) Rows() int {
	return b.rows
}

func (b *Board) Cols() int {
	return b.cols
}

func MakeBoard(rows, cols int) *Board {
	roundCols := ((cols-1)/8 + 1) * 8 // round up to multiple of 8 so it fits 64bit int
	return &Board{
		rows:   rows,
		cols:   cols,
		cells:  makeMatrix(rows, roundCols),
		temp:   makeMatrix(rows, roundCols),
		empty:  make([]byte, roundCols),
		colsum: make([]byte, roundCols),
	}
}

func makeMatrix(rows, cols int) [][]byte {
	all := make([]byte, rows*cols)
	c := make([][]byte, rows)
	for i := range c {
		c[i] = all[i*cols : (i+1)*cols]
	}
	return c
}

func zero(ps []byte) {
	for i := range ps {
		ps[i] = 0
	}
}
