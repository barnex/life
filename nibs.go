package life

import (
	"bytes"
	"fmt"
)

const (
	NibBits     = 4
	NibMask     = 0xF
	WordBits    = 64
	NibsPerWord = WordBits / NibBits
)

// Array of nibbles (4-bit words), packed in 64-bit ints
type Nibs []uint64

func (n Nibs) String() string {
	var buf bytes.Buffer
	for i := 0; i < n.nibs(); i++ {
		if i != 0 && i%NibsPerWord == 0 {
			fmt.Fprintf(&buf, " ")
		}
		fmt.Fprintf(&buf, "%x", n.get(i))
	}
	return buf.String()
}

func (n Nibs) get(i int) uint64 {
	w := i / NibsPerWord
	bitpos := uint(i % NibsPerWord)

	word := n[w]
	return getNib(word, bitpos)
}

func (n Nibs) set(i int, v uint64) {
	w := i / NibsPerWord
	nibpos := uint(i % NibsPerWord)
	n[w] = setNib(n[w], nibpos, uint64(v))
}

func getNib(w uint64, nibpos uint) uint64 {
	sh := nibpos * NibBits
	return (w >> sh) & NibMask
}

func setNib(w uint64, nibpos uint, v uint64) uint64 {
	if v > NibMask {
		panic(v)
	}
	sh := nibpos * NibBits
	w &^= (NibMask << sh)
	w |= (v << sh)
	return w
}

func (n Nibs) words() int {
	return len(n)
}

func (n Nibs) nibs() int {
	return n.words() * NibsPerWord
}

func makeNibs(n int) Nibs {
	if n%NibsPerWord != 0 {
		panic(n)
	}
	return make(Nibs, n/NibsPerWord)
}
