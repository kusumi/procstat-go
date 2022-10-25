//go:build !stdout

package main

import (
	"os"

	"github.com/rthornton128/goncurses"
)

type Screen struct {
	*goncurses.Window
}

var (
	stdscr        *Screen
	color_attr    goncurses.Char = goncurses.A_NORMAL
	standout_attr goncurses.Char = goncurses.A_NORMAL
)

const (
	KEY_ERR    = -1
	KEY_UP     = goncurses.KEY_UP
	KEY_DOWN   = goncurses.KEY_DOWN
	KEY_LEFT   = goncurses.KEY_LEFT
	KEY_RIGHT  = goncurses.KEY_RIGHT
	KEY_RESIZE = goncurses.KEY_RESIZE
)

func KEY_CTRL(x int) int {
	return x & 0x1F
}

func StringToColor(arg string) int16 {
	switch arg {
	case "black":
		return goncurses.C_BLACK
	case "red":
		return goncurses.C_RED
	case "green":
		return goncurses.C_GREEN
	case "yellow":
		return goncurses.C_YELLOW
	case "blue":
		return goncurses.C_BLUE
	case "magenta":
		return goncurses.C_MAGENTA
	case "cyan":
		return goncurses.C_CYAN
	case "white":
		return goncurses.C_WHITE
	}
	return -1 // default color
}

func InitScreen(fg, bg int16) error {
	scr, err := goncurses.Init()
	if err != nil {
		return err
	}
	stdscr = &Screen{scr}

	if err := goncurses.Cursor(0); err != nil {
		return err
	}

	goncurses.Echo(false)
	goncurses.CBreak(true)

	if err := stdscr.Keypad(true); err != nil {
		return err
	}

	stdscr.Timeout(200)
	ClearTerminal()

	if goncurses.HasColors() {
		if err := goncurses.StartColor(); err != nil {
			return err
		}
		if err := goncurses.UseDefaultColors(); err != nil {
			return err
		}
		if err := goncurses.InitPair(1, fg, bg); err != nil {
			return err
		}
		color_attr = goncurses.ColorPair(1)
	}

	if os.Getenv("TERM") == "screen" {
		standout_attr = goncurses.A_REVERSE
	} else {
		standout_attr = goncurses.A_STANDOUT
	}

	return nil
}

func CleanupScreen() {
	if err := goncurses.Cursor(1); err != nil {
		return
	}
	goncurses.End()
}

func ReadIncoming() int {
	return int(stdscr.GetChar())
}

func ClearTerminal() {
	GlobalLock()
	if err := stdscr.Clear(); err != nil {
		panic(err)
	}
	stdscr.Window.Refresh()
	GlobalUnlock()
}

func FlashTerminal() {
	goncurses.Flash()
}

func AllocScreen(ylen, xlen, ypos, xpos int) *Screen {
	GlobalLock()
	scr, err := goncurses.NewWindow(ylen, xlen, ypos, xpos)
	if err != nil {
		panic(err)
	}
	scr.ScrollOk(false)
	//scr.IdlOk(false) // XXX

	if err := scr.Keypad(true); err != nil {
		panic(err)
	}
	GlobalUnlock()

	return &Screen{scr}
}

func (this *Screen) Delete() {
	GlobalLock()
	_ = this.Window.Delete()
	GlobalUnlock()
}

func (this *Screen) Print(y, x int, standout bool, s string) {
	GlobalLock()
	attr := standout_attr
	if !standout {
		attr = goncurses.A_NORMAL
	}
	_ = this.Window.AttrOn(attr)
	this.Window.MovePrint(y, x, s)
	_ = this.Window.AttrOff(attr)
	GlobalUnlock()
}

func (this *Screen) Refresh() {
	GlobalLock()
	this.Window.Refresh()
	GlobalUnlock()
}

func (this *Screen) Erase() {
	GlobalLock()
	this.Window.Erase()
	GlobalUnlock()
}

func (this *Screen) Resize(ylen, xlen int) {
	GlobalLock()
	this.Window.Resize(ylen, xlen)
	GlobalUnlock()
}

func (this *Screen) Move(ypos, xpos int) {
	GlobalLock()
	this.Window.Move(ypos, xpos)
	GlobalUnlock()
}

func (this *Screen) Box() {
	GlobalLock()
	err := this.Window.Border(goncurses.ACS_VLINE, goncurses.ACS_VLINE,
		goncurses.ACS_HLINE, goncurses.ACS_HLINE,
		goncurses.ACS_ULCORNER, goncurses.ACS_URCORNER,
		goncurses.ACS_LLCORNER, goncurses.ACS_LRCORNER)
	if err != nil {
		panic(err)
	}
	GlobalUnlock()
}

func (this *Screen) Bkgd() {
	GlobalLock()
	if color_attr != goncurses.A_NORMAL {
		this.Window.SetBackground(color_attr | ' ')
	}
	GlobalUnlock()
}
