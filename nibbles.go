package life

import (
	"bytes"
	"fmt"
)

const (
	NibbleBits     = 4
	NibbleMask     = 0xF
	WordBits       = 64
	NibblesPerWord = WordBits / NibbleBits
)

// Array of nibbles (4-bit words), packed in 64-bit ints
type Nibbles []uint64

func (n Nibbles) String() string {
	var buf bytes.Buffer
	for i := 0; i < n.nibs(); i++ {
		if i != 0 && i%NibblesPerWord == 0 {
			fmt.Fprintf(&buf, " ")
		}
		fmt.Fprintf(&buf, "%x", n.get(i))
	}
	return buf.String()
}

func (n Nibbles) get(i int) uint64 {
	w := i / NibblesPerWord
	bitpos := uint(i % NibblesPerWord)

	word := n[w]
	return getNib(word, bitpos)
}

func (n Nibbles) set(i int, v uint64) {
	w := i / NibblesPerWord
	nibpos := uint(i % NibblesPerWord)
	n[w] = setNib(n[w], nibpos, uint64(v))
}

func getNib(w uint64, nibpos uint) uint64 {
	sh := nibpos * NibbleBits
	return (w >> sh) & NibbleMask
}

func setNib(w uint64, nibpos uint, v uint64) uint64 {
	if v > NibbleMask {
		panic(v)
	}
	sh := nibpos * NibbleBits
	w &^= (NibbleMask << sh)
	w |= (v << sh)
	return w
}

func (n Nibbles) words() int {
	return len(n)
}

func (n Nibbles) nibs() int {
	return n.words() * NibblesPerWord
}

func makeNibs(n int) Nibbles {
	if n%NibblesPerWord != 0 {
		panic(n)
	}
	return make(Nibbles, n/NibblesPerWord)
}

func makeMatrix(rows, cols int) []Nibbles {
	if cols%NibblesPerWord != 0 {
		panic(cols)
	}
	cols /= NibblesPerWord
	storage := make(Nibbles, rows*cols)
	c := make([]Nibbles, rows)
	for i := range c {
		c[i] = storage[i*cols : (i+1)*cols]
	}
	return c
}
