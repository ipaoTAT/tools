package tstat

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"strconv"
)

const (
	WINDOW_WIDTH_MIN  = 18
	WINDOW_HEIGHT_MIN = 10
)

const (
	CHAR_BLOCK = '#'
	CHAR_ALIGN = '.'
)

const (
	MAX_INT = int(^uint(0) >> 1)
	MIN_INT = int(^MAX_INT)
)

type StatWindow struct {
	// properties
	name       string
	title      string
	x, y, h, w int // input position
	// inner properties
	posX, posY, height, width int      // real position to control display
	point                     []int    // points
	buf                       [][]byte // display buffer
	// view
	*gocui.View
}

// function to new a statistic window
// name is the key, which must be unique in the same screen
// title is the display info on the top of window
// x, y show the position of the top-left point of the window area (include border)
// height, width is size of window area (include border)
func NewStatWindow(name, title string, x, y, height, width uint32) (*StatWindow, error) {
	// TODO check params
	if name == "" {
		return nil, errors.New("name is empty")
	}
	if height < WINDOW_HEIGHT_MIN {
		height = WINDOW_HEIGHT_MIN
	}
	if width < WINDOW_WIDTH_MIN {
		width = WINDOW_WIDTH_MIN
	}
	n := &StatWindow{}
	// fill properties
	n.name = name
	n.title = title
	n.x = int(x)
	n.y = int(y)
	n.h = int(height)
	n.w = int(width)
	n.posX, n.posY, n.height, n.width = n.x, n.y, n.h, n.w
	n.View = nil
	n.point = make([]int, 0)
	// init display buffer
	n.initBuffer()
	return n, nil
}

// function to init display buffer
func (w *StatWindow) initBuffer() {
	width, height := w.drawSize()
	w.buf = make([][]byte, height, height)
	for i := 0; i < len(w.buf); i++ {
		w.buf[i] = make([]byte, width, width)
	}
}

// function to clean display buffer
func (w *StatWindow) cleanBuffer() {
	for i := 0; i < len(w.buf); i++ {
		for j := 0; j < len(w.buf[i]); j++ {
			w.buf[i][j] = ' '
		}
	}
}

// function to fetch height and width of draw without broder
func (w *StatWindow) drawSize() (width, height int) {
	// no border
	return w.width - 2, w.height - 2
}

// function to render graph and output into display buffer
func (w *StatWindow) render() {
	width, height := w.drawSize()
	max, min := w.maxPoint(), w.minPoint()
	// compute unit for one line
	minAxisMark, step := CalculateAxisMark(height, max, min)
	markOfLine := func(l int) int {
		return (height-1-l)*step + minAxisMark
	}
	// clean buffer
	w.cleanBuffer()
	// render graph first
	for i := 1; i < width; i++ {
		valPos := len(w.point) - width + i
		if valPos >= 0 {
			for j := 0; j < height; j++ {
				// value is equal or greater than mark of this line
				if w.point[valPos] >= markOfLine(j) {
					w.buf[j][i] = CHAR_BLOCK
				}
			}
		}
	}
	// print y axis
	for i := 0; i < height; i++ {
		if (height-i-1)%5 == 0 {
			// draw y axis mark
			markStr := strconv.Itoa((markOfLine(i)))
			w.copyBytes(w.buf[i], 0, []byte(markStr))
			// draw align line
			for j := 0; j < len(w.buf[i]); j++ {
				if w.buf[i][j] == ' ' {
					w.buf[i][j] = CHAR_ALIGN
				}
			}
		}
	}
	// print currently data
	currData := 0
	if len(w.point) > 0 {
		currData = w.point[len(w.point)-1]
	}
	subWindowX, subWindowY := 5, 1
	message := fmt.Sprintf("Current: %d", currData)
	w.copyBytes(w.buf[subWindowY+1], subWindowX, []byte(message))
}

// function to refresh display
func (w *StatWindow) Refresh() {
	w.render()
	w.View.Clear()
	for _, bt := range w.buf {
		fmt.Fprint(w.View, string(bt))
	}
}

func (w *StatWindow) SetView(v *gocui.View) {
	v.Title = w.title
	v.Editable = false
	v.Wrap = true
	w.View = v
}

// function to get the max point for display
func (w *StatWindow) maxPoint() int {
	points := w.pointsToDisplay()
	if len(points) == 0 {
		return 0
	}
	max := MIN_INT
	for _, v := range points {
		if v > max {
			max = v
		}
	}
	return max
}

// function to get the min point for display
func (w *StatWindow) minPoint() int {
	points := w.pointsToDisplay()
	if len(points) == 0 {
		return 0
	}
	min := MAX_INT
	for _, v := range points {
		if v < min {
			min = v
		}
	}
	return min
}

// function to fetch points list to display
func (w *StatWindow) pointsToDisplay() []int {
	width, _ := w.drawSize()
	j := len(w.point) - width
	if j <= 0 {
		return w.point
	} else {
		return w.point[j:]
	}
}

// function to add new points to this statistic window
func (w *StatWindow) AddPoints(i ...int) {
	if len(i) == 0 {
		return
	}
	w.point = append(w.point, i...)
}

// function to clear all points
func (w *StatWindow) Reset() {
	w.point = make([]int, 0)
}

// function to set this statistic window in center mode
func (w *StatWindow) SetOnTop(g *gocui.Gui) {
	maxX, maxY := g.Size()
	w.posX, w.posY, w.height, w.width = 0, 0, maxY, maxX
	if w.posX+w.width > maxX {
		w.posX, w.width = 0, maxX
	}
	if w.posY+w.height > maxY {
		w.posY, w.height = 0, maxY
	}
	w.initBuffer()
	g.SetCurrentView(w.name)
	g.SetViewOnTop(w.name)
}

// function to switch focus off from this window
func (w *StatWindow) OnBlur(g *gocui.Gui) {
	w.posX, w.posY, w.height, w.width = w.x, w.y, w.h, w.w
	w.initBuffer()
}

// function to focus on this window
func (w *StatWindow) OnFocus(g *gocui.Gui) {
	maxX, maxY := g.Size()
	w.posX, w.posY, w.height, w.width = w.x-1, w.y-1, w.h+2, w.w+2
	if w.posX < 0 {
		w.posX, w.width = 0, w.w+1
	}
	if w.posY < 0 {
		w.posY, w.height = 0, w.h+1
	}
	if w.posX+w.width > maxX {
		w.width = maxX - w.posX
	}
	if w.posY+w.height > maxY {
		w.height = maxY - w.posY
	}
	w.initBuffer()
	g.SetCurrentView(w.name)
	g.SetViewOnTop(w.name)
}

// function to copy bytes of val to dst[start:]
func (w *StatWindow) copyBytes(dst []byte, start int, val []byte) {
	for i, bt := range val {
		dst[start+i] = bt
	}
}
