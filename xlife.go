//+build ignore

//This demonstration shows drawing entire generated images to an x window
//and calculating the speed of the generation and the drawing operations
//It should be noted that redrawing the entire image is inefficient, this demo
//was made to show me how fast I could draw to the windows with this method.
package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	//"github.com/BurntSushi/xgbutil/xevent"
	"image"
	"log"
	"os"
	"runtime/pprof"

	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/barnex/life"
)

const CellsPerWord = life.NibblesPerWord

var (
	Cols   = 1920 - 13
	Rows   = 1024
	Width  = life.DivUp(Cols, CellsPerWord) * CellsPerWord // image too wide to fit border
	Height = Rows
)

var board *life.Board

func main() {

	board = life.MakeBoard(Rows, Cols)
	life.SetRand(board, 1, 0.1)

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

	//So here i'm going to try to make a image..'
	img := xgraphics.New(X, image.Rect(0, 0, Width, Height))

	err = img.XSurfaceSet(win.Id)
	if err != nil {
		log.Printf("Error while setting window surface to image %d: %s\n",
			win, err)
	} else {
		log.Printf("Window %d surface set to image OK\n", win)
	}

	// I /think/ XDraw actually sends data to server?
	//img.XDraw()
	// I /think/ XPaint tells the server to paint image to window
	img.XPaint(win.Id)

	initPProf()

	//Start the routine that updates the window
	for i := 0; i < 128; i++ {
		updater(img, win)
	}

	pprof.StopCPUProfile()

	//This seems to start a main loop for listening to xevents
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

func updater(img *xgraphics.Image, win *xwindow.Window) {
	if gens == 32 {
		fps := float64(gens) / time.Since(start).Seconds()
		fmt.Println(fps, "FPS")
		start = time.Now()
		gens = 0
	}

	// render here
	board.Advance(1)
	board.Render(img.Pix)

	//hopefully using checked will block us from drawing again before x
	//draws although XDraw might block anyway, we can check for an error
	//here
	err := img.XDrawChecked()
	if err != nil {
		log.Println(err)
		return
	}
	//	img.XDraw()

	img.XPaint(win.Id)
	gens++
}
