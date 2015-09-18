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

func colSum(dst, up, me, down []byte) {
	for i := range dst {
		dst[i] = up[i] + me[i] + down[i]
	}
}

func (b *Board) advRow(r int, up, me, down []byte) {

	cs := b.colsum
	colSum(cs, up, me, down)

	cols := b.Cols()
	result := b.temp[r]

	// first col
	c := 0
	alive := me[c]
	neigh := cs[c] + cs[c+1]
	result[c] = nextLUT[(alive<<4)|neigh]

	// bulk cols
	for c := 1; c < cols-1; c++ {
		alive = me[c]
		neigh = cs[c-1] + cs[c] + cs[c+1]
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
