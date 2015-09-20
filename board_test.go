package life

import "fmt"

func ExampleCountNeigh() {
	rows, cols := 6, 3*16
	b := MakeBoard(rows, cols)

	BoardSet(b, 0, 0, [][]bool{
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X},
		{X, X, O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, X, X, O, O, O, X, O, O, O, O, O, O, O, X, O, X, O, O, O, O, O, O, O, O, O, O, O, O, X, X, X},
		{X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, X, X, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, X, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, X, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
	})

	count := makeMatrix(rows, cols)
	for r := range count {
		b.countNeigh(count[r], r)
	}
	for _, row := range count {
		fmt.Println(row)
	}

	// Output:
	//2210000000000111 0000000000001110 0000000000000011
	//3320000000000122 2221011100001222 1100000000001243
	//4420000000000122 3442011100001333 2210000000001343
	//2210000000000011 3553011212111333 2210000000001354
	//1100000000000000 1332000112111221 1110000000000122
	//0000000000000000 0111000112110000 0000000000000011

}

//func ExampleColSum() {
//	rows, cols := 6, 16
//	b := MakeBoard(rows, cols)
//
//	BoardSet(b, 0, 0, [][]bool{
//		{O, O, O, O, O},
//		{X, O, O, O, O},
//		{O, X, O, X, X},
//		{O, O, O, X, X},
//		{O, O, O, O, X},
//		{O, O, O, O, O},
//	})
//
//	rs := makeNibs(cols)
//	for r := 1; r < rows-1; r++ {
//		colSum(rs, b.cells[r-1], b.cells[r], b.cells[r+1])
//		for c := 0; c < rs.nibs(); c++ {
//			C := rs.get(c)
//			fmt.Print(C, ",")
//		}
//		fmt.Println()
//	}
//
//	// Output:
//	//1,1,0,1,1,0,0,0,0,0,0,0,0,0,0,0,
//	//1,1,0,2,2,0,0,0,0,0,0,0,0,0,0,0,
//	//0,1,0,2,3,0,0,0,0,0,0,0,0,0,0,0,
//	//0,0,0,1,2,0,0,0,0,0,0,0,0,0,0,0,
//}
