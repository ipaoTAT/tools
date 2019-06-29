package main

import (
	"github.com/jroimartin/gocui"
	"github.com/ipaoTAT/tools/golib/tstat"
)

var quit chan interface{}
var w1 *tstat.StatWindow
var g *tstat.StatScreen

func main() {
	var err error
	g, err = tstat.NewStatScreen()
	if err != nil {
		panic(err)
	}
	defer g.Close()
	x, y := g.Size()
	w1, err = tstat.NewStatWindow("v1", "Window XXXX", 0, 0, uint32(y), uint32(x))
	if err != nil {
		panic(err)
	}
	g.AddWindow(w1)
	quit = make(chan interface{})
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
	close(quit)
}