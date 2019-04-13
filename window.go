package main

import (
	"strconv"
	"sync"
)

type Window struct {
	frame  *Frame
	panel  *Panel
	buffer *Buffer
	offset int
	sig_ch chan bool
	mtx    sync.Mutex
}

func (this *Window) Init(ylen, xlen, ypos, xpos int) {
	this.frame = new(Frame)
	this.panel = new(Panel)
	this.frame.Init(ylen, xlen, ypos, xpos)
	this.panel.Init(ylen-2, xlen-2, ypos+1, xpos+1)
	this.frame.Refresh()
	this.panel.Refresh()

	this.offset = 0
	this.sig_ch = make(chan bool)
}

func (this *Window) Signal() {
	this.sig_ch <- true
}

func (this *Window) Lock() {
	this.mtx.Lock()
}

func (this *Window) Unlock() {
	this.mtx.Unlock()
}

func (this *Window) IsDead() bool {
	return this.buffer == nil || this.buffer.IsDead()
}

func (this *Window) Resize(ylen, xlen, ypos, xpos int) {
	this.Lock()
	this.frame.Resize(ylen, xlen, ypos, xpos)
	this.panel.Resize(ylen-2, xlen-2, ypos+1, xpos+1)
	this.Unlock()
}

func (this *Window) AttachBuffer(f string) {
	this.Lock()
	if this.buffer != nil {
		this.Unlock()
		return
	}
	this.frame.SetTitle(f)
	this.panel.SetTitle(f)

	this.buffer = new(Buffer)
	this.buffer.Init(f)
	dbgf("window=%p path=%s", this, this.buffer.f)
	this.Unlock()
}

func (this *Window) UpdateBuffer() {
	dbgf("window=%p path=%s", this, this.buffer.f)
	this.buffer.Update()
}

func (this *Window) Focus(t bool) {
	this.Lock()
	this.frame.SetFocus(t)
	this.panel.SetFocus(t)
	this.Unlock()
}

func (this *Window) GotoHead() {
	this.Lock()
	this.offset = 0
	this.Unlock()
}

func (this *Window) GotoTail() {
	this.Lock()
	this.offset = this.buffer.GetMaxLine()
	this.Unlock()
}

func (this *Window) GotoCurrent(d int) {
	this.Lock()
	x := this.offset + d
	if x < 0 {
		x = 0
	} else if x > this.buffer.GetMaxLine() {
		x = this.buffer.GetMaxLine()
	}
	this.offset = x
	this.Unlock()
}

func (this *Window) Repaint() {
	if this.IsDead() {
		return
	}

	l, err := this.buffer.ReadLines()
	if err != nil {
		this.panel.Erase()
		this.panel.Print(0, 0, false, err.Error())
		this.panel.Refresh()
		this.buffer.Clear()
		this.buffer.SaveLines([]string{})
		return
	}
	pl := this.buffer.prev
	y := 0

	this.Lock()
	offset := this.offset
	ylen := this.panel.ylen
	xlen := this.panel.xlen
	this.Unlock()

	this.buffer.BlockTillReady()
	this.panel.Erase()

	for i, s := range l {
		this.Lock()
		if y >= this.panel.ylen || offset != this.offset ||
			ylen != this.panel.ylen || xlen != this.panel.xlen {
			this.Unlock()
			break
		}
		this.Unlock()
		if i < offset {
			continue
		}
		standout := false
		if opt.blinkline && len(pl) > 0 {
			if i >= len(pl) || s != pl[i] {
				standout = true
			}
		}
		if opt.showlnum {
			s = strconv.Itoa(i) + " " + s
		}
		if !opt.foldline && len(s) > xlen {
			s = s[:xlen]
		}
		this.panel.Print(y, 0, standout, s)
		if !opt.foldline {
			y++
		} else {
			siz := len(s)
			y += siz / xlen
			if siz%xlen != 0 {
				y++
			}
		}
	}

	this.panel.Refresh()
	this.buffer.Clear()
	this.buffer.SaveLines(l)
	this.buffer.SignalBlocked()
}
