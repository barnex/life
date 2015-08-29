package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net/http"
)

func Fmt(b *Board) string {
	rows := b.Rows()
	cols := b.Cols()
	var buf bytes.Buffer
	fmt.Fprintln(&buf)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := b.Get(i, j)
			if v {
				fmt.Fprint(&buf, "X")
			} else {
				fmt.Fprint(&buf, ".")
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

func BoardSet(b *Board, roff, coff int, in [][]bool) {
	for r, row := range in {
		for c, v := range row {
			b.Set(r+roff, c+coff, v)
		}
	}
}

func SetRand(b *Board, seed int64, fill float64) {
	rand.Seed(seed)
	for r := 0; r < b.Rows(); r++ {
		for c := 0; c < b.Cols(); c++ {
			v := (rand.Float64() <= fill)
			b.Set(r, c, v)
		}
	}
}

func main() {
	const (
		X = true
		O = false
	)
	N := 256
	b := MakeBoard(N, N)
	SetRand(b, 0, 0.1)

	http.HandleFunc("/img", handleImg)
	http.HandleFunc("/", handleRoot)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	trash <- nil

	for {
		select {
		default:
			b.Advance(1)
		case img := <-request:
			rendered <- render(img, b)
			b.Advance(1) // make sure we advance at lease one step per rendering
		}
	}
}

var (
	trash    = make(chan *image.RGBA, 1)
	request  = make(chan *image.RGBA)
	rendered = make(chan *image.RGBA)
)

const JS = `
<html>
<head>
<script>
setInterval(function() {
    var el = document.getElementById('img');
    el.src = '/img?rand=' + Math.random();
}, 50);
</script>
</head>
<body>
	<img id="img" src="/img" alt="dead"/>
</body>
</html>
`

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, JS)
}

func handleImg(w http.ResponseWriter, r *http.Request) {
	request <- <-trash
	img := <-rendered
	png.Encode(w, img)
	trash <- img
}

var (
	Alive = color.Black
	Dead  = color.White
)

func render(img *image.RGBA, b *Board) *image.RGBA {
	rows, cols := b.Rows(), b.Cols()
	wantSize := image.Rect(0, 0, rows, cols)
	if img == nil || img.Bounds() != wantSize {
		img = image.NewRGBA(image.Rect(0, 0, rows, cols))
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if b.Get(r, c) {
				img.Set(r, c, Alive)
			} else {
				img.Set(r, c, Dead)
			}
		}
	}

	return img
}
