package life

import (
	"runtime"
	"unsafe"
)

type Board struct {
	rows, cols int
	cells      [][]byte // current cells
	temp       [][]byte // buffer for next-gen cells
	empty      []byte   // empty cell row used at borders
	work, done chan int
}

// Advance the state given number of steps
func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

// advance one step
func (b *Board) advance() {
	// do rows in parallel
	for r := 0; r < b.rows; r++ {
		b.work <- r
	}
	// wait for result
	for r := 0; r < b.rows; r++ {
		<-b.done
	}
	// swap: temp becomes current cells
	b.cells, b.temp = b.temp, b.cells
}

// view byte array as 64-bit int array,
// so we can add 8 pairs of bytes in one instruction.
func as64(bytes []byte) []uint64 {
	if len(bytes)%8 != 0 {
		panic("as64")
	}
	return (*((*[1 << 31]uint64)(unsafe.Pointer(&bytes[0]))))[:len(bytes)/8]
}

// partial sums towards number of neighbors:
// dst = up + me + down, per column
// Arrays must have multiple of 8 size so we can add SIMD-style
func colSum(dst, up, me, down []byte) {
	dst64 := as64(dst)
	up64 := as64(up)
	me64 := as64(me)
	down64 := as64(down)

	for i := range dst64 {
		dst64[i] = up64[i] + me64[i] + down64[i]
	}
}

// advance row r to the next state,
// freely using cs as a buffer.
func (b *Board) advRow(r int, cs []byte) {

	// get adjacent rows without going out of bounds
	up := b.empty
	if r > 0 {
		up = b.cells[r-1]
	}
	me := b.cells[r]
	down := b.empty
	if r < b.rows-1 {
		down = b.cells[r+1]
	}

	cols := b.Cols()
	result := b.temp[r]

	colSum(cs, up, me, down)

	// partial column sums left, centered and right of current cell
	var prevCS, currCS, nextCS byte

	// first column is special
	c := 0
	alive := me[c]

	prevCS = 0
	currCS = cs[c]
	nextCS = cs[c+1]

	neigh := prevCS + currCS + nextCS
	result[c] = nextLUT[(alive<<4)|neigh]

	// bulk columns don't have borders
	for c := 1; c < cols-1; c++ {
		alive = me[c]

		prevCS = currCS
		currCS = nextCS
		nextCS = cs[c+1]

		neigh = prevCS + currCS + nextCS
		result[c] = nextLUT[(alive<<4)|neigh]
	}

	// last column is special
	c = cols - 1
	alive = me[c]
	neigh = cs[c-1] + cs[c]
	result[c] = nextLUT[(alive<<4)|neigh]
}

// look-up table for next state,
// indexed by (alive<<4)|(neigh+alive)
var nextLUT [32]byte

// set-up nextLUT
func init() {
	for _, alive := range []byte{0, 1} {
		for neigh := byte(0); neigh <= 8; neigh++ {
			idx := (alive << 4) | neigh
			nextLUT[idx] = nextState(alive, neigh-alive) // self is included in neigh
		}
	}
}

// next cell state for current alive state and number of neighbors
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
	b := &Board{
		rows:  rows,
		cols:  cols,
		cells: makeMatrix(rows, roundCols),
		temp:  makeMatrix(rows, roundCols),
		empty: make([]byte, roundCols),
		work:  make(chan int, rows),
		done:  make(chan int, rows),
	}

	// start parallel workers:
	// TODO: give more than one row per worker
	// so that small boards run efficiently as well.
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			colsum := make([]byte, roundCols)
			for {
				r := <-b.work
				b.advRow(r, colsum)
				b.done <- 1
			}
		}()
	}

	return b
}

func makeMatrix(rows, cols int) [][]byte {
	all := make([]byte, rows*cols)
	c := make([][]byte, rows)
	for i := range c {
		c[i] = all[i*cols : (i+1)*cols]
	}
	return c
}
