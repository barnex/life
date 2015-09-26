//+build ignore

package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/barnex/life"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	N := 512
	board := life.MakeBoard(N, N)
	life.SetRand(board, 0, 0.1)
	img := render(nil, board)

	// Now convert it into an X image.
	ximg := xgraphics.NewConvert(X, img)

	// Now show it in a new window.
	// We set the window title and tell the program to quit gracefully when
	// the window is closed.
	// There is also a convenience method, XShow, that requires no parameters.
	ximg.XShowExtra("The Go Gopher!", true)

	// If we don't block, the program will end and the window will disappear.
	// We could use a 'select{}' here, but xevent.Main will emit errors if
	// something went wrong, so use that instead.
	xevent.Main(X)
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
				img.Set(r, c, color.Black)
			} else {
				img.Set(r, c, color.White)
			}
		}
	}

	return img
}
