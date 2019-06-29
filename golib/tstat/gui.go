package tstat

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

var ErrWindowNotFound = errors.New("cannot find window")

const (
	HELP_VIEW_NAME = "help"
	HELP_MESSAGE   = `
	HOT KEYS:
		<Tab>  : switch window
		<Enter>: set current window into center mode
		?      : set this help page into center mode
		q      : quit center mode
		<Ctl+c>: exit program

	Type 'q' to quit this help
	`
)

// Screen of statistics
type StatScreen struct {
	*gocui.Gui
	windows     map[string]*StatWindow
	windowNames []string
}

// function to new a statistic screen, if any unexpected, return error
func NewStatScreen() (*StatScreen, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorRed
	ret := &StatScreen{g, make(map[string]*StatWindow), nil}
	g.SetManagerFunc(ret.layout)
	ret.mainModeBindingKeys()
	return ret, nil
}

// init key bindings for main mode
func (g *StatScreen) mainModeBindingKeys() {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, g.quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, g.switchWindow); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, g.setWindowOnTop); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, g.unsetWindowOnTop); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", '?', gocui.ModNone, g.enterHelpMode); err != nil {
		panic(err)
	}
}

// init key bindings for help mode
func (g *StatScreen) helpModeBindingKeys() {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, g.quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, g.quitHelpMode); err != nil {
		panic(err)
	}
}

// clear all key bindings
func (g *StatScreen) clearBindingKeys() {
	g.DeleteKeybindings("")
}

// function to management layout, render all windows
func (statG *StatScreen) layout(g *gocui.Gui) error {
	for _, name := range statG.windowNames {
		w, ok := statG.windows[name]
		if !ok || w == nil {
			return ErrWindowNotFound
		}
		if v, err := g.SetView(w.name, w.posX, w.posY, w.posX+w.width-1, w.posY+w.height-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			w.SetView(v)
			// set the first window as the current window and set it on focus
			if statG.CurrentView() == nil {
				w.OnFocus(g)
			} else {
				statG.SetViewOnTop(statG.CurrentView().Name())
			}
		}
		w.Refresh()
	}
	return nil
}

// function to add a statistic window into screen
func (g *StatScreen) AddWindow(w *StatWindow) {
	g.windows[w.name] = w
	g.windowNames = append(g.windowNames, w.name)
}

// function to get statistic window from screen with name,
// if window doesn't exist, return error 'ErrWindowNotFound'
func (g *StatScreen) GetWindow(name string) (*StatWindow, error) {
	w, ok := g.windows[name]
	if !ok || w == nil {
		return nil, ErrWindowNotFound
	}
	return w, nil
}

// add points into window with name 'windowName' in this screen
// if window doesn't exist, return error 'ErrWindowNotFound'
func (g *StatScreen) AddPoints(windowName string, points ...int) error {
	window, ok := g.windows[windowName]
	if !ok {
		return ErrWindowNotFound
	}
	window.AddPoints(points...)
	g.Refresh()
	return nil
}

// force refresh screen display
func (g *StatScreen) Refresh() {
	g.Update(func(*gocui.Gui) error { return nil })
}

// key-binding functions
func (statG *StatScreen) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func (statG *StatScreen) switchWindow(g *gocui.Gui, v *gocui.View) error {
	// find current window
	currName := statG.CurrentView().Name()
	curr, ok := statG.windows[currName]
	if !ok {
		return ErrWindowNotFound
	}
	curr.OnBlur(g)
	nextWindowName := currName
	for i, n := range statG.windowNames {
		if n == currName {
			pos := (i + 1) % len(statG.windowNames)
			nextWindowName = statG.windowNames[pos]
			break
		}
	}
	w, ok := statG.windows[nextWindowName]
	if !ok {
		return ErrWindowNotFound
	}
	w.OnFocus(g)
	return nil
}
func (statG *StatScreen) setWindowOnTop(g *gocui.Gui, v *gocui.View) error {
	w, ok := statG.windows[statG.CurrentView().Name()]
	if !ok {
		return ErrWindowNotFound
	}
	w.SetOnTop(g)
	return nil
}
func (statG *StatScreen) unsetWindowOnTop(g *gocui.Gui, v *gocui.View) error {
	w, ok := statG.windows[statG.CurrentView().Name()]
	if !ok {
		return ErrWindowNotFound
	}
	w.OnFocus(g)
	return nil
}
func (statG *StatScreen) enterHelpMode(g *gocui.Gui, v *gocui.View) error {
	if !statG.inHelpMode() {
		maxX, maxY := g.Size()
		if v, err := g.SetView(HELP_VIEW_NAME, 5, 5, maxX-5, maxY-5); err != nil {
			// new help view
			if err != gocui.ErrUnknownView {
				return err
			}
			fmt.Fprint(v, HELP_MESSAGE)
			g.SetViewOnTop(HELP_VIEW_NAME)
			statG.clearBindingKeys()
			statG.helpModeBindingKeys()
		}
	}
	return nil
}
func (statG *StatScreen) quitHelpMode(g *gocui.Gui, v *gocui.View) error {
	if statG.inHelpMode() {
		g.DeleteView(HELP_VIEW_NAME)
		statG.clearBindingKeys()
		statG.mainModeBindingKeys()
	}
	return nil
}

// if this screen in help mode
func (g *StatScreen) inHelpMode() bool {
	v, _ := g.View(HELP_VIEW_NAME)
	return v != nil
}
