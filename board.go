// Package life provides a SIMD-accelerated implementation of
// Conway's game of life.
package life

// Board stores cells states and provides a method for advancing to the next generation.
type Board struct {
	rows, cols int
	Cells      []Nibbles // current cells
	temp       []Nibbles // buffer for next-gen cells
	empty      Nibbles   // empty cell row used at borders
	//work, done  chan int // for multi-threading
}

// Advance the state given number of steps
func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		for r := range b.Cells {
			b.advanceRow(r)
		}
		// swap: temp becomes current cells
		b.Cells, b.temp = b.temp, b.Cells
	}
}

// countNeigh counts the neighbors of all cells on row r
// and stores the result in dst.
func (b *Board) countNeigh(dst Nibbles, r int) {

	prevRow, currRow, nextRow := b.adjacentRows(r)

	// pipeline with adjacent per-column words
	var prev, curr, next uint64
	next = prevRow[0] + currRow[0] + nextRow[0] // prime the pipeline with first column sum

	// offset by one for easy retrieval of next element
	prevRow = prevRow[1:]
	currRow = currRow[1:]
	nextRow = nextRow[1:]

	// bulk cells have no boundary problems
	i := 0
	for ; i < len(dst)-1; i++ {
		prev = curr
		curr = next
		next = prevRow[i] + currRow[i] + nextRow[i]

		shr := (curr << NibbleBits) | (prev >> (WordBits - NibbleBits))
		shl := (curr >> NibbleBits) | (next << (WordBits - NibbleBits))

		dst[i] = shr + curr + shl
	}

	// last word has no next.
	prev = curr
	curr = next
	next = 0

	shr := (curr << NibbleBits) | (prev >> (WordBits - NibbleBits))
	shl := (curr >> NibbleBits) | (next & NibbleMask)

	dst[i] = shr + curr + shl

}

// advance row r to the next state,
func (b *Board) advanceRow(r int) {

	row := b.Cells[r]
	dst := b.temp[r]
	b.countNeigh(dst, r) // abuse dst to temporarily store neighbor count

	for w, alive := range row {
		ngbr := dst[w]
		keys := ngbr | (alive << 3)

		idx := uint16(keys)
		next := LUT4[idx]
		out := next
		keys >>= 16

		idx = uint16(keys)
		next = LUT4[idx]
		out |= next << 16
		keys >>= 16

		idx = uint16(keys)
		next = LUT4[idx]
		out |= next << 32
		keys >>= 16

		idx = uint16(keys)
		next = LUT4[idx]
		out |= next << 48

		dst[w] = out
	}

	// truncate last row
	// TODO: speed-up a bit
	for c := b.cols; c < dst.nibs(); c++ {
		dst.Set(c, 0)
	}
}

// get rows adjacent to r, without going out of bounds.
// returns row r-1, r, r+1, replacing out-of-bound rows
// by a row of zeros.
func (b *Board) adjacentRows(r int) (prev, curr, next Nibbles) {
	prev = b.empty
	if r > 0 {
		prev = b.Cells[r-1]
	}
	curr = b.Cells[r]
	next = b.empty
	if r < b.rows-1 {
		next = b.Cells[r+1]
	}
	return prev, curr, next
}

// Massive look-up table for the next states of four cells at a time.
// The key for one cell is (alive<<3)|(neigh+alive), with
// neigh+alive the number of neighbors including the cell itself.
var LUT4 [16 * 16 * 16 * 16]uint64

// set-up LUT4
func init() {
	// table for one cell
	var lut [16]uint64
	for _, alive := range []uint64{0, 1} {
		for neigh := uint64(0); neigh <= 8; neigh++ {
			idx := (alive << 3) | neigh
			lut[idx] = nextState(alive, neigh-alive) // self is included in neigh
		}
	}

	// table for two cells at a time
	var LUT2 [16 * 16]uint64
	for i1, v1 := range lut {
		for i2, v2 := range lut {
			I := (i1 << NibbleBits) | i2
			V := (v1 << NibbleBits) | v2
			LUT2[I] = V
		}
	}

	// table for four cells at a time
	for i1, v1 := range LUT2 {
		for i2, v2 := range LUT2 {
			I := (i1 << 8) | i2
			V := (v1 << 8) | v2
			LUT4[I] = V
		}
	}
}

// next cell state for current alive state and number of neighbors
func nextState(alive uint64, neighbors uint64) uint64 {
	if alive == 1 && neighbors == 2 || neighbors == 3 {
		return 1
	}
	return 0
}

func (b *Board) Set(r, c int, v bool) {
	if v {
		b.Cells[r].Set(c, 1)
	} else {
		b.Cells[r].Set(c, 0)
	}
}

func (b *Board) get(r, c int) uint64 {
	if r < 0 || c < 0 || r >= b.Rows() || c >= b.Cols() {
		return 0
	}
	return b.Cells[r].Get(c)
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

// Integer division, rounded up.
func DivUp(x, y int) int {
	return ((x-1)/y + 1)
}

func MakeBoard(rows, cols int) *Board {
	roundCols := DivUp(cols, NibblesPerWord) * NibblesPerWord // round up number of columns to fit 64-bit
	b := &Board{
		rows:  rows,
		cols:  cols,
		Cells: makeMatrix(rows, roundCols),
		temp:  makeMatrix(rows, roundCols),
		empty: makeNibs(roundCols),
	}

	// start parallel workers:
	// TODO: give more than one row per worker
	// so that small boards run efficiently as well.
	//for i := 0; i < runtime.NumCPU(); i++ {
	//	go func() {
	//		colsum := make([]byte, roundCols)
	//		for {
	//			r := <-b.work
	//			b.advRow(r, colsum)
	//			b.done <- 1
	//		}
	//	}()
	//}

	return b
}

//func (b *Board) stepParallel() {
//	// do rows in parallel
//	for r := 0; r < b.rows; r++ {
//		b.work <- r
//	}
//	// wait for result
//	for r := 0; r < b.rows; r++ {
//		<-b.done
//	}
//}
