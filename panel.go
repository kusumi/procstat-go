package main

type Panel struct {
	ylen, xlen, ypos, xpos int
	scr                    *Screen
}

type Frame struct {
	Panel
	title string
	focus bool
}

// Panel
func (this *Panel) Init(ylen, xlen, ypos, xpos int) {
	this.scr = AllocScreen(ylen, xlen, ypos, xpos)
	this.setYX(ylen, xlen, ypos, xpos)
	this.scr.Bkgd()
}

func (this *Panel) setYX(ylen, xlen, ypos, xpos int) {
	this.ylen = ylen
	this.xlen = xlen
	this.ypos = ypos
	this.xpos = xpos
}

func (this *Panel) Refresh() {
	this.scr.Refresh()
}

func (this *Panel) Erase() {
	this.scr.Erase()
}

func (this *Panel) Resize(ylen, xlen, ypos, xpos int) {
	this.doResize(ylen, xlen, ypos, xpos)
	this.Refresh()
}

func (this *Panel) doResize(ylen, xlen, ypos, xpos int) {
	this.setYX(ylen, xlen, ypos, xpos)
	this.scr.Resize(ylen, xlen)
	this.scr.Move(ypos, xpos)
}

func (this *Panel) SetTitle(s string) {
}

func (this *Panel) SetFocus(t bool) {
}

func (this *Panel) Print(y int, x int, standout bool, s string) {
	this.scr.Print(y, x, standout, s)
}

// Frame
func (this *Frame) Init(ylen, xlen, ypos, xpos int) {
	this.Panel.Init(ylen, xlen, ypos, xpos)
	this.scr.Box()
}

func (this *Frame) Refresh() {
	this.Panel.Refresh()
}

func (this *Frame) Erase() {
	this.Panel.Erase()
}

func (this *Frame) Resize(ylen, xlen, ypos, xpos int) {
	this.doResize(ylen, xlen, ypos, xpos)
	this.scr.Box()
	this.PrintTitle()
}

func (this *Frame) SetTitle(s string) {
	this.title = s
	this.PrintTitle()
}

func (this *Frame) SetFocus(t bool) {
	this.focus = t
	this.PrintTitle()
}

func (this *Frame) PrintTitle() {
	this.Panel.Print(0, 1, this.focus, this.title)
	this.Refresh()
}
