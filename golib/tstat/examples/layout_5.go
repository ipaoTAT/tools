package main

import (
	"github.com/jroimartin/gocui"
	"github.com/ipaoTAT/tools/golib/tstat"
	//"math/rand"
	"time"
	"math/rand"
)

var quit chan interface{}
var w1, w2, w3, w4, w5 *tstat.StatWindow
var g *tstat.StatScreen

func main() {
	var err error
	g, err = tstat.NewStatScreen()
	if err != nil {
		panic(err)
	}
	defer g.Close()
	x, y := g.Size()
	w1, err = tstat.NewStatWindow("v1", "Statistic of port-1", 0, 0, uint32(y/2), uint32(x/3))
	if err != nil {
		panic(err)
	}
	w2, err = tstat.NewStatWindow("v2", "Statistic of port-2", 0, uint32(y/2), uint32(y/2), uint32(x/3))
	if err != nil {
		panic(err)
	}
	w3, err = tstat.NewStatWindow("v3", "Statistic of port-3", uint32(x/3), 0, uint32(y/2), uint32(x/3))
	if err != nil {
		panic(err)
	}
	w4, err = tstat.NewStatWindow("v4", "Statistic of port-4", uint32(x/3), uint32(y/2), uint32(y/2), uint32(x/3))
	if err != nil {
		panic(err)
	}
	w5, err = tstat.NewStatWindow("v5", "Statistic of port-5", uint32(x*2/3), 0, uint32(y), uint32(x/3))
	if err != nil {
		panic(err)
	}
	g.AddWindow(w1)
	g.AddWindow(w2)
	g.AddWindow(w3)
	g.AddWindow(w4)
	g.AddWindow(w5)
	quit = make(chan interface{})
	go addPoint()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
	close(quit)
}

func addPoint() {
	tk := time.NewTicker(time.Millisecond * 200)
	var i = 0
	for {
		select {
		case <-quit:
			return
		case <-tk.C:
			i = rand.Intn(10000)
			g.AddPoints("v1", i)
			i = rand.Intn(10000)
			g.AddPoints("v2", i)
			i = rand.Intn(10000)
			g.AddPoints("v3", i)
			i = rand.Intn(10000)
			g.AddPoints("v4", i)
			i = rand.Intn(10000)
			g.AddPoints("v5", i)
		}
	}
}
