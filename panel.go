package main

import (
	"github.com/rthornton128/goncurses"
)

type Panel struct {
	ylen, xlen, ypos, xpos int
	scr                    *goncurses.Window
}

type Frame struct {
	Panel
	title string
	focus bool
}

// Panel
func (this *Panel) Init(ylen, xlen, ypos, xpos int) {
	this.scr = AllocScreen(ylen, xlen, ypos, xpos)
	this.ylen = ylen
	this.xlen = xlen
	this.ypos = ypos
	this.xpos = xpos
	BkgdScreen(this.scr)
}

func (this *Panel) Refresh() {
	RefreshScreen(this.scr)
}

func (this *Panel) Erase() {
	EraseScreen(this.scr)
}

func (this *Panel) Resize(ylen, xlen, ypos, xpos int) {
	this.doResize(ylen, xlen, ypos, xpos)
	this.Refresh()
}

func (this *Panel) doResize(ylen, xlen, ypos, xpos int) {
	this.ylen = ylen
	this.xlen = xlen
	this.ypos = ypos
	this.xpos = xpos
	ResizeScreen(this.scr, ylen, xlen)
	MoveScreen(this.scr, ypos, xpos)
}

func (this *Panel) SetTitle(s string) {
}

func (this *Panel) SetFocus(t bool) {
}

func (this *Panel) Print(y int, x int, standout bool, s string) {
	PrintScreen(this.scr, y, x, standout, s)
}

// Frame
func (this *Frame) Init(ylen, xlen, ypos, xpos int) {
	this.Panel.Init(ylen, xlen, ypos, xpos)
	BoxScreen(this.scr)
}

func (this *Frame) Refresh() {
	this.Panel.Refresh()
}

func (this *Frame) Erase() {
	this.Panel.Erase()
}

func (this *Frame) Resize(ylen, xlen, ypos, xpos int) {
	this.doResize(ylen, xlen, ypos, xpos)
	BoxScreen(this.scr)
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
