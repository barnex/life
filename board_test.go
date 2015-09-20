package life

import "fmt"

func ExampleCountNeigh() {
	rows, cols := 6, 16
	b := MakeBoard(rows, cols)

	BoardSet(b, 0, 0, [][]bool{
		{O, O, O, O, O, O},
		{X, O, O, O, O, O},
		{O, X, O, X, X, O},
		{O, O, O, X, X, O},
		{O, O, O, O, X, O},
		{O, O, O, O, O, O},
	})

	count := makeMatrix(rows, cols)
	for r := range count {
		b.countNeigh(count[r], r)
	}
	for _, row := range count {
		fmt.Println(row)
	}

	// Output:
	//1100000000000000
	//2222210000000000
	//1135530000000000
	//0003320000000000
	//0000100000000000
}

func ExampleColSum() {
	rows, cols := 6, 16
	b := MakeBoard(rows, cols)

	BoardSet(b, 0, 0, [][]bool{
		{O, O, O, O, O},
		{X, O, O, O, O},
		{O, X, O, X, X},
		{O, O, O, X, X},
		{O, O, O, O, X},
		{O, O, O, O, O},
	})

	rs := makeNibs(cols)
	for r := 1; r < rows-1; r++ {
		colSum(rs, b.cells[r-1], b.cells[r], b.cells[r+1])
		for c := 0; c < rs.nibs(); c++ {
			C := rs.get(c)
			fmt.Print(C, ",")
		}
		fmt.Println()
	}

	// Output:
	//1,1,0,1,1,0,0,0,0,0,0,0,0,0,0,0,
	//1,1,0,2,2,0,0,0,0,0,0,0,0,0,0,0,
	//0,1,0,2,3,0,0,0,0,0,0,0,0,0,0,0,
	//0,0,0,1,2,0,0,0,0,0,0,0,0,0,0,0,
}
