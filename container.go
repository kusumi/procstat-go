package main

import (
	"fmt"
	"os"

	"github.com/rthornton128/goncurses"
)

type Container struct {
	v, bv []*Window
	curw  *Window
}

func (this *Container) Init(args []string) {
	this.BuildWindow()
	for i, f := range args {
		if _, err := os.Stat(f); err == nil {
			if i < len(this.v) {
				this.v[i].AttachBuffer(f)
				this.bv = append(this.bv, this.v[i])
			}
		}
	}

	Assert(len(this.v) > 0)
	this.curw = this.v[0]
	this.curw.Focus(true)
}

func (this *Container) GotoNextWindow() {
	this.curw.Focus(false)
	for i, w := range this.bv {
		if w == this.curw {
			if i == len(this.bv)-1 {
				this.curw = this.bv[0]
			} else {
				this.curw = this.bv[i+1]
			}
			this.curw.Focus(true)
			return
		}
	}
	Assert(len(this.bv) == 0)
}

func (this *Container) GotoPrevWindow() {
	this.curw.Focus(false)
	for i, w := range this.bv {
		if w == this.curw {
			if i == 0 {
				this.curw = this.bv[len(this.bv)-1]
			} else {
				this.curw = this.bv[i-1]
			}
			this.curw.Focus(true)
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

func KBD_CTRL(x int) int {
	return x & 0x1F
}

func (this *Container) ParseEvent(x int) int {
	switch x {
	case -1:
		dbg("KEY_ERR")
	case goncurses.KEY_RESIZE, KBD_CTRL('l'):
		ClearTerminal()
		this.BuildWindow()
	case goncurses.KEY_LEFT, 'h':
		this.GotoPrevWindow()
	case goncurses.KEY_RIGHT, 'l':
		this.GotoNextWindow()
	case '0':
		this.curw.GotoHead()
		this.curw.Signal()
	case '$':
		this.curw.GotoTail()
		this.curw.Signal()
	case goncurses.KEY_UP, 'k':
		this.curw.GotoCurrent(-1)
		this.curw.Signal()
	case goncurses.KEY_DOWN, 'j':
		this.curw.GotoCurrent(1)
		this.curw.Signal()
	case KBD_CTRL('B'):
		this.curw.GotoCurrent(-GetTerminalLines())
		this.curw.Signal()
	case KBD_CTRL('U'):
		this.curw.GotoCurrent(-GetTerminalLines() / 2)
		this.curw.Signal()
	case KBD_CTRL('F'):
		this.curw.GotoCurrent(GetTerminalLines())
		this.curw.Signal()
	case KBD_CTRL('D'):
		this.curw.GotoCurrent(GetTerminalLines() / 2)
		this.curw.Signal()
	}

	return 0
}
