package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type container struct {
	v, bv []*window
	cw    *window
}

func (this *container) init(args []string, watch *watch) {
	this.buildWindow()
	for i, f := range args {
		if st, err := os.Stat(f); err == nil && st.Mode().IsRegular() {
			if i < len(this.v) {
				this.v[i].attachBuffer(f)
				this.bv = append(this.bv, this.v[i])

				abs, _ := filepath.Abs(f)
				if err := watch.watcher.Add(abs); err == nil {
					watch.fmap[abs] = this.v[i]
				} else {
					dbg("Failed to watch", abs, err)
				}
			}
		}
	}

	assert(len(this.v) > 0)
	this.cw = this.v[0]
	this.cw.focus(true)
}

func (this *container) cleanup() {
	for _, w := range this.v {
		w.cleanup()
	}
}

func (this *container) gotoNextWindow() {
	this.cw.focus(false)
	for i, w := range this.bv {
		if w == this.cw {
			if i == len(this.bv)-1 {
				this.cw = this.bv[0]
			} else {
				this.cw = this.bv[i+1]
			}
			this.cw.focus(true)
			return
		}
	}
	if len(this.bv) > 0 {
		this.cw = this.bv[0]
		this.cw.focus(true)
	}
}

func (this *container) gotoPrevWindow() {
	this.cw.focus(false)
	for i, w := range this.bv {
		if w == this.cw {
			if i == 0 {
				this.cw = this.bv[len(this.bv)-1]
			} else {
				this.cw = this.bv[i-1]
			}
			this.cw.focus(true)
			return
		}
	}
	if len(this.bv) > 0 {
		this.cw = this.bv[len(this.bv)-1]
		this.cw.focus(true)
	}
}

func (this *container) buildWindow() {
	if !opt.rotatecol {
		this.buildWindowXY()
	} else {
		this.buildWindowYX()
	}
}

func (this *container) buildWindowXY() {
	seq := 0
	xx := getTerminalCols()
	yy := getTerminalLines()
	x := len(opt.layout)
	xq := xx / x
	xr := xx % x

	for i := 0; i < x; i++ {
		xpos := xq * i
		xlen := xq
		if i == x-1 {
			xlen += xr
		}
		y := opt.layout[i]
		if y == -1 {
			y = 1 // ignore invalid
		}
		yq := yy / y
		yr := yy % y

		for j := 0; j < y; j++ {
			ypos := yq * j
			ylen := yq
			if j == y-1 {
				ylen += yr
			}
			this.allocWindow(seq, ylen, xlen, ypos, xpos)
			seq++
		}
	}
}

func (this *container) buildWindowYX() {
	seq := 0
	yy := getTerminalLines()
	xx := getTerminalCols()
	y := len(opt.layout)
	yq := yy / y
	yr := yy % y

	for i := 0; i < y; i++ {
		ypos := yq * i
		ylen := yq
		if i == y-1 {
			ylen += yr
		}
		x := opt.layout[i]
		if x == -1 {
			x = 1 // ignore invalid
		}
		xq := xx / x
		xr := xx % x

		for j := 0; j < x; j++ {
			xpos := xq * j
			xlen := xq
			if j == x-1 {
				xlen += xr
			}
			this.allocWindow(seq, ylen, xlen, ypos, xpos)
			seq++
		}
	}
}

func (this *container) allocWindow(seq, ylen, xlen, ypos, xpos int) {
	s := fmt.Sprintf("#%d %d %d %d %d", seq, ylen, xlen, ypos, xpos)
	if len(this.v) > seq {
		w := this.v[seq]
		w.resize(ylen, xlen, ypos, xpos)
		dbgf("window=%p resize %s", w, s)
	} else {
		w := &window{}
		w.init(ylen, xlen, ypos, xpos)
		this.v = append(this.v, w)
		assert(this.v[seq] == w)
		dbgf("window=%p alloc %s", w, s)
	}
}

func (this *container) parseEvent(x int) int {
	switch x {
	case KEY_ERR:
		dbg("KEY_ERR")
	case KEY_RESIZE, keyCtrl('l'):
		clearTerminal()
		this.buildWindow()
	case KEY_LEFT, 'h':
		this.gotoPrevWindow()
	case KEY_RIGHT, 'l':
		this.gotoNextWindow()
	case '0':
		this.cw.gotoHead()
		this.cw.signal()
	case '$':
		this.cw.gotoTail()
		this.cw.signal()
	case KEY_UP, 'k':
		this.cw.gotoCurrent(-1)
		this.cw.signal()
	case KEY_DOWN, 'j':
		this.cw.gotoCurrent(1)
		this.cw.signal()
	case keyCtrl('B'):
		this.cw.gotoCurrent(-getTerminalLines())
		this.cw.signal()
	case keyCtrl('U'):
		this.cw.gotoCurrent(-getTerminalLines() / 2)
		this.cw.signal()
	case keyCtrl('F'):
		this.cw.gotoCurrent(getTerminalLines())
		this.cw.signal()
	case keyCtrl('D'):
		this.cw.gotoCurrent(getTerminalLines() / 2)
		this.cw.signal()
	}

	return 0
}
