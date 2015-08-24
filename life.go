package main

type Board struct {
	cells [][]bool
}

func (b *Board) Neighbors(r, c int) int {
	count := 0
	for rr := r - 1; rr <= r+1; rr++ {
		for cc := c - 1; cc <= c+1; cc++ {
			if b.Get(rr, cc) {
				count++
			}
		}
	}
	// do not count self
	if b.Get(r, c) {
		count--
	}
	return count
}

func (b *Board) Set(r, c int, v bool) {
	r, c = b.wrap(r, c)
	b.cells[r][c] = v
}

func (b *Board) Get(r, c int) bool {
	r, c = b.wrap(r, c)
	return b.cells[r][c]
}
func (b *Board) Rows() int {
	return len(b.cells)
}

func (b *Board) Cols() int {
	return len(b.cells[0])
}

func MakeBoard(rows, cols int) *Board {
	all := make([]bool, rows*cols)
	c := make([][]bool, rows)
	for i := range c {
		c[i] = all[i*cols : (i+1)*cols]
	}
	return &Board{cells: c}
}

func (b *Board) wrap(r, c int) (int, int) {
	r = (r + b.Rows()) % b.Rows()
	c = (c + b.Cols()) % b.Cols()
	return r, c
}

func main() {

}
