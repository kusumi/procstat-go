package main

type panel struct {
	ylen, xlen, ypos, xpos int
	scr                    screen
}

type frame struct {
	panel
	title string
	focus bool
}

// panel
func (this *panel) init(ylen, xlen, ypos, xpos int) {
	this.scr = allocScreen(ylen, xlen, ypos, xpos)
	this.ylen = ylen
	this.xlen = xlen
	this.ypos = ypos
	this.xpos = xpos
	this.scr.bkgd()
}

func (this *panel) cleanup() {
	this.scr.delete()
}

func (this *panel) refresh() {
	this.scr.refresh()
}

func (this *panel) erase() {
	this.scr.erase()
}

func (this *panel) resize(ylen, xlen, ypos, xpos int) {
	this._resize(ylen, xlen, ypos, xpos)
	this.refresh()
}

func (this *panel) _resize(ylen, xlen, ypos, xpos int) {
	// XXX goncurses resize only works with new window allocation
	this.init(ylen, xlen, ypos, xpos)
	this.scr.resize(ylen, xlen)
	this.scr.move(ypos, xpos)
}

func (this *panel) setTitle(s string) {
}

func (this *panel) setFocus(t bool) {
}

func (this *panel) print(y int, x int, standout bool, s string) {
	this.scr.print(y, x, standout, s)
}

// frame
func (this *frame) init(ylen, xlen, ypos, xpos int) {
	this.panel.init(ylen, xlen, ypos, xpos)
	this.scr.box()
}

func (this *frame) cleanup() {
	this.scr.delete()
}

func (this *frame) refresh() {
	this.panel.refresh()
}

func (this *frame) resize(ylen, xlen, ypos, xpos int) {
	this._resize(ylen, xlen, ypos, xpos)
	this.scr.box()
	this.printTitle()
}

func (this *frame) setTitle(s string) {
	this.title = s
	this.printTitle()
}

func (this *frame) setFocus(t bool) {
	this.focus = t
	this.printTitle()
}

func (this *frame) printTitle() {
	this.panel.print(0, 1, this.focus, this.title)
	this.refresh()
}
