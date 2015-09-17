package life

type Board struct {
	cells                     [][]byte // cell life states
	temp                      [][]byte // next state, swap with cells after advance
	edgerow                   []byte   // row with zeros used at edges
	sumPrev, sumCurr, sumNext []byte   // partial sums for upper, me, down row
}

// Advance advances the state by a number of steps.
func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

func (b *Board) advance() {
	rows := b.Rows()

	sumPrev := b.sumPrev
	sumCurr := b.sumCurr
	sumNext := b.sumNext

	r := 0
	zero(sumPrev)
	b.rowSum(r, sumCurr)
	b.rowSum(r+1, sumNext)
	b.advRow(r, b.edgerow, sumCurr, sumNext)

	for r = 1; r < rows-1; r++ {
		sumPrev = sumCurr
		sumCurr = sumNext
		b.rowSum(r+1, sumNext)
		b.advRow(r, sumPrev, sumCurr, sumNext)
	}

	r = rows - 1
	sumPrev = sumCurr
	sumCurr = sumNext
	b.advRow(r, sumPrev, sumCurr, b.edgerow)

	b.cells, b.temp = b.temp, b.cells
}

func (b *Board) rowSum(r int, result []byte) {
	row := b.cells[r]
	last := len(row) - 1

	var prev byte = 0
	var curr byte = row[0]
	var next byte = row[1]

	result[0] = curr + next
	for i := 1; i < last; i++ {
		prev = curr
		curr = next
		next = row[i+1]
		result[i] = prev + curr + next
	}

	result[last] = row[last-1] + row[last]
}

func (b *Board) advRow(r int, sumPrev, sumCurr, sumNext []byte) {
	rowIn := b.cells[r]
	rowOut := b.temp[r]
	for c, alive := range rowIn {
		neigh := sumPrev[c] + sumCurr[c] + sumNext[c] - alive
		rowOut[c] = nextState(alive, neigh)
	}
}

func (b *Board) advSlow(r, c int, up, me, down []byte) {
	alive := b.cells[r][c]
	neigh := b.neighbors(r, c)
	b.temp[r][c] = nextState(alive, neigh)
}

func nextState(alive byte, neighbors byte) byte {
	if alive == 1 && neighbors == 2 || neighbors == 3 {
		return 1
	}
	return 0
}

func (b *Board) Neighbors(r, c int) int {
	return int(b.neighbors(r, c))
}

func (b *Board) neighbors(r, c int) byte {
	cL := c - 1
	cR := c + 1

	count := b.get(r, cL)
	count += b.get(r, cR)

	r--
	count += b.get(r, cL)
	count += b.get(r, c)
	count += b.get(r, cR)

	r += 2
	count += b.get(r, cL)
	count += b.get(r, c)
	count += b.get(r, cR)

	return count
}

func (b *Board) innerNeigh(up, me, do []byte, c int) byte {
	cL := c - 1
	cR := c + 1

	cUp := up[cL] + up[c] + up[cR]
	cMe := me[cL] + me[cR]
	cDo := do[cL] + do[c] + do[cR]

	return cUp + cMe + cDo
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
	return len(b.cells)
}

func (b *Board) Cols() int {
	return len(b.cells[0])
}

func MakeBoard(rows, cols int) *Board {
	return &Board{
		cells:   makeMatrix(rows, cols),
		temp:    makeMatrix(rows, cols),
		edgerow: make([]byte, cols),
		sumPrev: make([]byte, cols),
		sumCurr: make([]byte, cols),
		sumNext: make([]byte, cols),
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
