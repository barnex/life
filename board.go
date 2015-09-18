package life

type Board struct {
	cells  [][]byte // current cells
	temp   [][]byte // buffer for next-gen cells
	empty  []byte   // empty cell row used at borders
	colsum []byte   //buffer for vertical sums by 3
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

func (b *Board) advRow(r int, up, me, down []byte) {
	cols := b.Cols()
	b.advSlow(r, 0, up, me, down)
	b.advInner(r, up, me, down)
	b.advSlow(r, cols-1, up, me, down)
}

func colSum(dst, up, me, down []byte) {
	for i := range dst {
		dst[i] = up[i] + me[i] + down[i]
	}
}

func (b *Board) advInner(r int, up, me, down []byte) {

	cs := b.colsum
	colSum(cs, up, me, down)

	cols := b.Cols()
	result := b.temp[r]
	for c := 1; c < cols-1; c++ {
		alive := me[c]
		cL := c - 1
		cR := c + 1
		neigh := up[cL] + up[c] + up[cR] +
			me[cL] + me[cR] +
			down[cL] + down[c] + down[cR]
		result[c] = nextState(alive, neigh)
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
		cells:  makeMatrix(rows, cols),
		temp:   makeMatrix(rows, cols),
		empty:  make([]byte, cols),
		colsum: make([]byte, cols),
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
