package main

type Board struct {
	cells [][]byte
	temp  [][]byte
}

func (b *Board) Advance(steps int) {
	for i := 0; i < steps; i++ {
		b.advance()
	}
}

func (b *Board) advance() {
	for i := range b.cells {
		for j := range b.cells[i] {
			neigh := b.neighbors(i, j)
			alive := b.cells[i][j]
			b.temp[i][j] = nextState(alive, neigh)
		}
	}
	b.cells, b.temp = b.temp, b.cells
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

	r--
	count := b.get(r, cL)
	count += b.get(r, c)
	count += b.get(r, cR)

	r++
	count += b.get(r, cL)
	count += b.get(r, cR)

	r++
	count += b.get(r, cL)
	count += b.get(r, c)
	count += b.get(r, cR)

	return count
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
		cells: makeMatrix(rows, cols),
		temp:  makeMatrix(rows, cols),
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
