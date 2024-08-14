//go:build !stdout

package main

import (
	"os"

	"github.com/rthornton128/goncurses"
)

type screen struct {
	win *goncurses.Window
}

var (
	stdscr       *goncurses.Window
	colorAttr    goncurses.Char = goncurses.A_NORMAL
	standoutAttr goncurses.Char = goncurses.A_NORMAL
)

const (
	KEY_ERR    = -1
	KEY_UP     = goncurses.KEY_UP
	KEY_DOWN   = goncurses.KEY_DOWN
	KEY_LEFT   = goncurses.KEY_LEFT
	KEY_RIGHT  = goncurses.KEY_RIGHT
	KEY_RESIZE = goncurses.KEY_RESIZE
)

func keyCtrl(x int) int {
	return x & 0x1F
}

func stringToColor(arg string) int16 {
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

func initScreen(fg, bg int16) error {
	scr, err := goncurses.Init()
	if err != nil {
		return err
	}
	stdscr = scr

	if err := goncurses.Cursor(0); err != nil {
		return err
	}

	goncurses.Echo(false)
	goncurses.CBreak(true)

	if err := stdscr.Keypad(true); err != nil {
		return err
	}

	stdscr.Timeout(200)
	clearTerminal()

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
		colorAttr = goncurses.ColorPair(1)
	}

	if os.Getenv("TERM") == "screen" {
		standoutAttr = goncurses.A_REVERSE
	} else {
		standoutAttr = goncurses.A_STANDOUT
	}

	return nil
}

func cleanupScreen() {
	if err := goncurses.Cursor(1); err != nil {
		return
	}
	goncurses.End()
}

func readIncoming() int {
	return int(stdscr.GetChar())
}

func clearTerminal() {
	globalLock()
	if err := stdscr.Clear(); err != nil {
		panic(err)
	}
	stdscr.Refresh()
	globalUnlock()
}

func flashTerminal() {
	goncurses.Flash()
}

func allocScreen(ylen, xlen, ypos, xpos int) screen {
	globalLock()
	scr, err := goncurses.NewWindow(ylen, xlen, ypos, xpos)
	if err != nil {
		panic(err)
	}
	scr.ScrollOk(false)
	//scr.IdlOk(false) // XXX

	if err := scr.Keypad(true); err != nil {
		panic(err)
	}
	globalUnlock()

	return screen{scr}
}

func (this *screen) delete() {
	globalLock()
	_ = this.win.Delete()
	globalUnlock()
}

func (this *screen) print(y, x int, standout bool, s string) {
	globalLock()
	attr := standoutAttr
	if !standout {
		attr = goncurses.A_NORMAL
	}
	_ = this.win.AttrOn(attr)
	this.win.MovePrint(y, x, s)
	_ = this.win.AttrOff(attr)
	globalUnlock()
}

func (this *screen) refresh() {
	globalLock()
	this.win.Refresh()
	globalUnlock()
}

func (this *screen) erase() {
	globalLock()
	this.win.Erase()
	globalUnlock()
}

func (this *screen) resize(ylen, xlen int) {
	globalLock()
	this.win.Resize(ylen, xlen)
	globalUnlock()
}

func (this *screen) move(ypos, xpos int) {
	globalLock()
	this.win.Move(ypos, xpos)
	globalUnlock()
}

func (this *screen) box() {
	globalLock()
	err := this.win.Border(goncurses.ACS_VLINE, goncurses.ACS_VLINE,
		goncurses.ACS_HLINE, goncurses.ACS_HLINE,
		goncurses.ACS_ULCORNER, goncurses.ACS_URCORNER,
		goncurses.ACS_LLCORNER, goncurses.ACS_LRCORNER)
	if err != nil {
		panic(err)
	}
	globalUnlock()
}

func (this *screen) bkgd() {
	globalLock()
	if colorAttr != goncurses.A_NORMAL {
		this.win.SetBackground(colorAttr | ' ')
	}
	globalUnlock()
}
