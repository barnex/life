package life

import (
	"unsafe"
)

const (
	White uint64 = 0x00FFFFFF
	Black uint64 = 0x00000000
)

func (b *Board) Render(pix []uint8) {

	Words := len(b.cells[0])
	Width := Words * NibblesPerWord
	Height := b.Rows()
	pixels := (*(*[1 << 31]uint64)(unsafe.Pointer(&pix[0])))[:Width*Height]

	i := 0
	for _, row := range b.cells {
		for _, word := range row {
			pixels[i] = lut2[byte(word)]
			i++
			pixels[i] = lut2[byte(word>>8)]
			i++
			pixels[i] = lut2[byte(word>>16)]
			i++
			pixels[i] = lut2[byte(word>>24)]
			i++
			pixels[i] = lut2[byte(word>>32)]
			i++
			pixels[i] = lut2[byte(word>>40)]
			i++
			pixels[i] = lut2[byte(word>>48)]
			i++
			pixels[i] = lut2[byte(word>>56)]
			i++
		}
	}
}

var lut2 [16 * 16]uint64

func init() {

	var lut [16]uint64
	lut[1] = White

	for k1 := range lut {
		for k2 := range lut {
			v1 := lut[k1]
			v2 := lut[k2]
			K := k1<<NibbleBits | k2
			V := v1<<32 | v2
			lut2[K] = V
		}
	}
}
