//+build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/barnex/life"
)

func main() {

	X, err := xgb.NewConn()
	if err != nil {
		fmt.Println(err)
		return
	}

	wid, _ := xproto.NewWindowId(X)
	screen := xproto.Setup(X).DefaultScreen(X)
	xproto.CreateWindow(X, screen.RootDepth, wid, screen.Root,
		0, 0, 500, 500, 0,
		xproto.WindowClassInputOutput, screen.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{ // values must be in the order defined by the protocol
			0xffffffff,
			xproto.EventMaskStructureNotify |
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease})

	xproto.MapWindow(X, wid)

	time.Sleep(10 * time.Second)

}

func render(img *image.RGBA, b *life.Board) *image.RGBA {
	rows, cols := b.Rows(), b.Cols()
	wantSize := image.Rect(0, 0, rows, cols)
	if img == nil || img.Bounds() != wantSize {
		img = image.NewRGBA(image.Rect(0, 0, rows, cols))
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if b.Get(r, c) {
				img.Set(r, c, color.White)
			} else {
				img.Set(r, c, color.Black)
			}
		}
	}

	return img
}
