package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Container struct {
	v, bv []*Window
	cw    *Window
}

func (this *Container) Init(args []string, watch *Watch) {
	this.BuildWindow()
	for i, f := range args {
		if st, err := os.Stat(f); err == nil && st.Mode().IsRegular() {
			if i < len(this.v) {
				this.v[i].AttachBuffer(f)
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

	Assert(len(this.v) > 0)
	this.cw = this.v[0]
	this.cw.Focus(true)
}

func (this *Container) GotoNextWindow() {
	this.cw.Focus(false)
	for i, w := range this.bv {
		if w == this.cw {
			if i == len(this.bv)-1 {
				this.cw = this.bv[0]
			} else {
				this.cw = this.bv[i+1]
			}
			this.cw.Focus(true)
			return
		}
	}
	Assert(len(this.bv) == 0)
}

func (this *Container) GotoPrevWindow() {
	this.cw.Focus(false)
	for i, w := range this.bv {
		if w == this.cw {
			if i == 0 {
				this.cw = this.bv[len(this.bv)-1]
			} else {
				this.cw = this.bv[i-1]
			}
			this.cw.Focus(true)
			return
		}
	}
	Assert(len(this.bv) == 0)
}

func (this *Container) BuildWindow() {
	if !opt.rotatecol {
		this.BuildWindowXY()
	} else {
		this.BuildWindowYX()
	}
}

func (this *Container) BuildWindowXY() {
	seq := 0
	xx := GetTerminalCols()
	yy := GetTerminalLines()
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
			this.AllocWindow(seq, ylen, xlen, ypos, xpos)
			seq++
		}
	}
}

func (this *Container) BuildWindowYX() {
	seq := 0
	yy := GetTerminalLines()
	xx := GetTerminalCols()
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
			this.AllocWindow(seq, ylen, xlen, ypos, xpos)
			seq++
		}
	}
}

func (this *Container) AllocWindow(seq, ylen, xlen, ypos, xpos int) {
	s := fmt.Sprintf("#%d %d %d %d %d", seq, ylen, xlen, ypos, xpos)
	if len(this.v) > seq {
		w := this.v[seq]
		w.Resize(ylen, xlen, ypos, xpos)
		dbgf("window=%p resize %s", w, s)
	} else {
		w := &Window{}
		w.Init(ylen, xlen, ypos, xpos)
		this.v = append(this.v, w)
		Assert(this.v[seq] == w)
		dbgf("window=%p alloc %s", w, s)
	}
}

func KEY_CTRL(x int) int {
	return x & 0x1F
}

func (this *Container) ParseEvent(x int) int {
	switch x {
	case KEY_ERR:
		dbg("KEY_ERR")
	case KEY_RESIZE, KEY_CTRL('l'):
		ClearTerminal()
		this.BuildWindow()
	case KEY_LEFT, 'h':
		this.GotoPrevWindow()
	case KEY_RIGHT, 'l':
		this.GotoNextWindow()
	case '0':
		this.cw.GotoHead()
		this.cw.Signal()
	case '$':
		this.cw.GotoTail()
		this.cw.Signal()
	case KEY_UP, 'k':
		this.cw.GotoCurrent(-1)
		this.cw.Signal()
	case KEY_DOWN, 'j':
		this.cw.GotoCurrent(1)
		this.cw.Signal()
	case KEY_CTRL('B'):
		this.cw.GotoCurrent(-GetTerminalLines())
		this.cw.Signal()
	case KEY_CTRL('U'):
		this.cw.GotoCurrent(-GetTerminalLines() / 2)
		this.cw.Signal()
	case KEY_CTRL('F'):
		this.cw.GotoCurrent(GetTerminalLines())
		this.cw.Signal()
	case KEY_CTRL('D'):
		this.cw.GotoCurrent(GetTerminalLines() / 2)
		this.cw.Signal()
	}

	return 0
}
