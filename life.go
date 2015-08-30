package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
)

func main() {
	const (
		X = true
		O = false
	)
	N := 128
	b := MakeBoard(N, N)

	BoardSet(b, 0, 0, [][]bool{
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, X, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O},
		{O, X, X, O, O, O, O, O, O, O, O, X, O, O, O, O, O, X, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, X, X, O, O, O, O, O, O, O, O, X, O, O, O, X, O, X, X, O, O, O, O, X, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, X, O, O, O, O, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
	})

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
	if (typeof img.naturalWidth !== "undefined" && img.naturalWidth === 0){
   		el.src = '/img?rand=' + Math.random();
	}
}, 100);
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
