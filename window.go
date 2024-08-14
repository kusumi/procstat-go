package main

import (
	"strconv"
)

type window struct {
	frame
	panel
	buffer
	offset int
	sigCh  chan int
	mtxCh  chan int
}

func (this *window) init(ylen, xlen, ypos, xpos int) {
	this.frame.init(ylen, xlen, ypos, xpos)
	this.panel.init(ylen-2, xlen-2, ypos+1, xpos+1)
	this.frame.refresh()
	this.panel.refresh()

	this.offset = 0
	this.sigCh = make(chan int)
	this.mtxCh = make(chan int, 1)
	this.mtxCh <- 1
}

func (this *window) cleanup() {
	this.frame.cleanup()
	this.panel.cleanup()
}

func (this *window) signal() {
	this.sigCh <- 1
}

func (this *window) lock() {
	<-this.mtxCh
}

func (this *window) unlock() {
	this.mtxCh <- 1
}

func (this *window) isDead() bool {
	return this.buffer.isDead()
}

func (this *window) resize(ylen, xlen, ypos, xpos int) {
	this.lock()
	this.frame.resize(ylen, xlen, ypos, xpos)
	this.panel.resize(ylen-2, xlen-2, ypos+1, xpos+1)
	this.unlock()
}

func (this *window) attachBuffer(f string) {
	this.lock()
	if this.buffer.fd != nil {
		this.unlock()
		return
	}
	this.frame.setTitle(f)
	this.panel.setTitle(f)

	this.buffer.init(f)
	dbgf("window=%p path=%s", this, this.buffer.f)
	this.unlock()
}

func (this *window) updateBuffer() {
	dbgf("window=%p path=%s", this, this.buffer.f)
	this.buffer.update()
}

func (this *window) focus(t bool) {
	this.lock()
	this.frame.setFocus(t)
	this.panel.setFocus(t)
	this.unlock()
}

func (this *window) gotoHead() {
	this.lock()
	this.offset = 0
	this.unlock()
}

func (this *window) gotoTail() {
	this.lock()
	this.offset = this.buffer.getMaxLine()
	this.unlock()
}

func (this *window) gotoCurrent(d int) {
	this.lock()
	x := this.offset + d
	if x < 0 {
		x = 0
	} else if x > this.buffer.getMaxLine() {
		x = this.buffer.getMaxLine()
	}
	this.offset = x
	this.unlock()
}

func (this *window) repaint() {
	if this.isDead() {
		return
	}

	l, err := this.buffer.readLines()
	if err != nil {
		this.panel.erase()
		this.panel.print(0, 0, false, err.Error())
		this.panel.refresh()
		this.buffer.clear()
		this.buffer.saveLines([]string{})
		return
	}
	pl := this.buffer.prev
	y := 0

	this.lock()
	offset := this.offset
	ylen := this.panel.ylen
	xlen := this.panel.xlen
	this.unlock()

	this.buffer.blockTillReady()
	this.panel.erase()

	for i, s := range l {
		this.lock()
		if y >= this.panel.ylen || offset != this.offset ||
			ylen != this.panel.ylen || xlen != this.panel.xlen {
			this.unlock()
			break
		}
		this.unlock()
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
		this.panel.print(y, 0, standout, s)
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

	this.panel.refresh()
	this.buffer.clear()
	this.buffer.saveLines(l)
	this.buffer.signalBlocked()
}
