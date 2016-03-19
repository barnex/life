/*
Command xlife runs game of life and shows the state in an X window.
It uses modified code from https://github.com/BurntSushi/xgbutil.
*/
package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/barnex/life"
)

const CellsPerWord = life.NibblesPerWord

var (
	Cols   = 1920
	Rows   = 1080
	Width  = life.DivUp(Cols, CellsPerWord) * CellsPerWord // image too wide to fit border
	Height = Rows
)

var board *life.Board

func main() {

	board = life.MakeBoard(Rows, Cols)
	//life.SetRand(board, 1, 0.02)

	const (
		O = false
		I = true
	)

	glider := [][]bool{
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, I, O, O, I, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, I, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, I, O, O, O, I, O, O, O, O, O, I, O, O},
		{O, O, O, O, O, O, O, I, I, I, I, O, O, O, O, I, O, O, O, I},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, I, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, I, O, O, O, O, I},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, I, I, I, I, I, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
	}

	for i := 0; i < Rows-100; i++ {
		for j := 0; j < Cols-100; j++ {
			if rand.Float64() < 0.00003 {
				life.BoardSet(board, i, j, glider)
			}
		}
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Println(err)
		return
	}

	//Initialize the keybind package
	keybind.Initialize(X)

	//Create a window
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatalf("Could not generate a new window X id: %s", err)
	}
	win.Create(X.RootWin(), 0, 0, Width, Height, xproto.CwBackPixel, 0x606060ff)

	// Listen for Key{Press,Release} events.
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease)

	// Map the window. This is what makes it on the screen
	win.Map()

	//Make a ...callback... for the events and connect
	//xevent.KeyPressFun(
	//	func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
	//		modStr := keybind.ModifierString(e.State)
	//		keyStr := keybind.LookupString(X, e.State, e.Detail)
	//		if len(modStr) > 0 {
	//			log.Printf("Key: %s-%s\n", modStr, keyStr)
	//		} else {
	//			log.Println("Key:", keyStr)
	//		}

	//		if keybind.KeyMatch(X, "Escape", e.State, e.Detail) {
	//			if e.State&xproto.ModMaskControl > 0 {
	//				log.Println("Control-Escape detected. Quitting...")
	//				xevent.Quit(X)
	//			}
	//		}
	//	}).Connect(X, win.Id)

	img := xgraphics.New(X, image.Rect(0, 0, Width, Height))

	time.Sleep(2 * time.Second)

	err = img.XSurfaceSet(win.Id)
	if err != nil {
		log.Printf("Error while setting window surface to image %d: %s\n",
			win, err)
	} else {
		log.Printf("Window %d surface set to image OK\n", win)
	}

	img.XPaint(win.Id)

	//initPProf()

	//Start the routine that updates the window
	for i := 0; i < 128; i += 0 {
		updater(img, win)
	}

	//pprof.StopCPUProfile()

	//xevent.Main(X)
}

func initPProf() {
	fname := "cpu.pprof"
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("writing CPU profile to", fname)
}

var (
	start = time.Now()
	gens  int
)

const GensPerStep = 15

func updater(img *xgraphics.Image, win *xwindow.Window) {
	if gens == 32 {
		fps := float64(gens) / time.Since(start).Seconds()
		fmt.Println(fps, "FPS", (fps*GensPerStep*float64(Rows*Cols))/1e6, "Mc/s")
		start = time.Now()
		gens = 0
	}

	// render here
	board.Advance(GensPerStep)
	board.Render(img.Pix)

	err := img.XDrawChecked()
	if err != nil {
		log.Println(err)
		return
	}
	//	img.XDraw()

	img.XPaint(win.Id)
	gens++
}
