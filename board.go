package life

type Board struct {
	rows, cols            int
	cells                 []Nibs // current cells
	temp                  []Nibs // buffer for next-gen cells
	empty                 Nibs   // empty cell row used at borders
	work, done            chan int
	serialCS, serialNeigh Nibs // temp hack for serial tuning
}

// Advance the state given number of steps
func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

// advance one step
func (b *Board) advance() {
	//b.stepParallel()
	b.stepSerial()

	// swap: temp becomes current cells
	b.cells, b.temp = b.temp, b.cells
}

func (b *Board) stepSerial() {
	for r := range b.cells {
		b.advRow(r, b.serialCS, b.serialNeigh)
	}
}

func (b *Board) stepParallel() {
	// do rows in parallel
	for r := 0; r < b.rows; r++ {
		b.work <- r
	}
	// wait for result
	for r := 0; r < b.rows; r++ {
		<-b.done
	}
}

// dst[i] = a[i] + b[i] + c[i].
func colSum(dst, a, b, c Nibs) {
	for i := range dst {
		dst[i] = a[i] + b[i] + c[i]
	}
}

func (b *Board) countNeigh(dst, cs Nibs, r int) {

	prevRow, currRow, nextRow := b.adjacentRows(r)
	colSum(cs, prevRow, currRow, nextRow)

	// pipeline adjacent colSum words
	var prev, curr, next uint64
	next = cs[0] // prime the pipeline

	// bulk
	i := 0
	for ; i < len(cs)-1; i++ {
		prev = curr
		curr = next
		next = cs[i+1]

		shr := (curr << NibBits) | (prev >> (WordBits - NibBits))
		shl := (curr >> NibBits) | (next << (WordBits - NibBits))

		dst[i] = shr + curr + shl
	}

	// last
	prev = curr
	curr = next
	next = 0

	shr := (curr << NibBits) | (prev >> (WordBits - NibBits))
	shl := (curr >> NibBits) | (next & NibMask)

	dst[i] = shr + curr + shl

}

// advance row r to the next state,
// freely using neigh as a buffer.
func (b *Board) advRow(r int, buf1, buf2 Nibs) {

	row := b.cells[r]
	dst := b.temp[r]
	cs := buf1
	neigh := buf2

	b.countNeigh(neigh, cs, r)

	c := 0
	for w := 0; w < dst.words(); w++ {

		alive := row[w]
		ngbr := neigh[w]
		keys := ngbr | (alive << 3)

		var out uint64
		for n := uint(0); n < NibsPerWord; n++ {

			idx := keys & NibMask
			next := nextLUT[idx]
			out = setNib(out, n, next)

			keys >>= NibBits
			c++
		}
		dst[w] = out

	}

	// truncate last row
	for c := b.cols; c < dst.nibs(); c++ {
		dst.set(c, 0)
	}
}

// get rows adjacent to r, without going out of bounds.
// returns row r-1, r, r+1, replacing out-of-bound rows
// by a row of zeros.
func (b *Board) adjacentRows(r int) (prev, curr, next Nibs) {
	prev = b.empty
	if r > 0 {
		prev = b.cells[r-1]
	}
	curr = b.cells[r]
	next = b.empty
	if r < b.rows-1 {
		next = b.cells[r+1]
	}
	return prev, curr, next
}

// look-up table for next state,
// indexed by (alive<<3)|(neigh+alive)
var nextLUT [16]uint64

// set-up nextLUT
func init() {
	for _, alive := range []uint64{0, 1} {
		for neigh := uint64(0); neigh <= 8; neigh++ {
			idx := (alive << 3) | neigh
			nextLUT[idx] = nextState(alive, neigh-alive) // self is included in neigh
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
		b.cells[r].set(c, 1)
	} else {
		b.cells[r].set(c, 0)
	}
}

func (b *Board) get(r, c int) uint64 {
	if r < 0 || c < 0 || r >= b.Rows() || c >= b.Cols() {
		return 0
	}
	return b.cells[r].get(c)
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
	roundCols := ((cols-1)/NibsPerWord + 1) * NibsPerWord // round up to multiple of 8 so it fits 64bit int
	b := &Board{
		rows:        rows,
		cols:        cols,
		cells:       makeMatrix(rows, roundCols),
		temp:        makeMatrix(rows, roundCols),
		empty:       makeNibs(roundCols),
		work:        make(chan int, rows),
		done:        make(chan int, rows),
		serialCS:    makeNibs(roundCols),
		serialNeigh: makeNibs(roundCols),
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

// TODO: contiguous!!
func makeMatrix(rows, cols int) []Nibs {
	c := make([]Nibs, rows)
	for i := range c {
		c[i] = makeNibs(cols)
	}
	return c
}
