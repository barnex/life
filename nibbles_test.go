package life

import "fmt"

func ExampleNibbles_Set() {
	n := makeNibs(32)
	for i := 0; i < n.nibs(); i++ {
		n.Set(i, uint64(i%10))
	}
	fmt.Println(n)

	//Output:
	//0123456789012345 6789012345678901
}

func ExampleSetNib() {
	w := uint64(0)
	for i := uint(0); i < NibblesPerWord; i++ {
		w = setNib(w, i, uint64(i))
		fmt.Printf("%016x\n", w)
	}

	// Output:
	//0000000000000000
	//0000000000000010
	//0000000000000210
	//0000000000003210
	//0000000000043210
	//0000000000543210
	//0000000006543210
	//0000000076543210
	//0000000876543210
	//0000009876543210
	//00000a9876543210
	//0000ba9876543210
	//000cba9876543210
	//00dcba9876543210
	//0edcba9876543210
	//fedcba9876543210
}
