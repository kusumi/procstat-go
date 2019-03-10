// +build !stdout

package main

import (
	"os"

	"github.com/rthornton128/goncurses"
)

var (
	stdscr        *goncurses.Window
	color_attr    goncurses.Char = goncurses.A_NORMAL
	standout_attr goncurses.Char = goncurses.A_NORMAL
)

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

func InitScreen(fg, bg int16) {
	scr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	stdscr = scr

	if err := goncurses.Cursor(0); err != nil {
		panic(err)
	}

	goncurses.Echo(false)
	goncurses.CBreak(true)

	if err := stdscr.Keypad(true); err != nil {
		panic(err)
	}

	stdscr.Timeout(200)
	ClearTerminal()

	if goncurses.HasColors() == true {
		if err := goncurses.StartColor(); err != nil {
			panic(err)
		}
		if err := goncurses.UseDefaultColors(); err != nil {
			panic(err)
		}
		if err := goncurses.InitPair(1, fg, bg); err != nil {
			panic(err)
		}
		color_attr = goncurses.ColorPair(1)
	}

	if os.Getenv("TERM") == "screen" {
		standout_attr = goncurses.A_REVERSE
	} else {
		standout_attr = goncurses.A_STANDOUT
	}
}

func CleanupScreen() {
	if err := goncurses.Cursor(1); err != nil {
		panic(err)
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
	stdscr.Refresh()
	GlobalUnlock()
}

func FlashTerminal() {
	goncurses.Flash()
}

func AllocScreen(ylen, xlen, ypos, xpos int) *goncurses.Window {
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

	return scr
}

func DeleteScreen(scr *goncurses.Window) {
	GlobalLock()
	scr.Delete()
	GlobalUnlock()
}

func PrintScreen(scr *goncurses.Window, y int, x int, standout bool, s string) {
	GlobalLock()
	attr := standout_attr
	if !standout {
		attr = goncurses.A_NORMAL
	}
	scr.AttrOn(attr)
	scr.MovePrint(y, x, s)
	scr.AttrOff(attr)
	GlobalUnlock()
}

func RefreshScreen(scr *goncurses.Window) {
	GlobalLock()
	scr.Refresh()
	GlobalUnlock()
}

func EraseScreen(scr *goncurses.Window) {
	GlobalLock()
	scr.Erase()
	GlobalUnlock()
}

func ResizeScreen(scr *goncurses.Window, ylen int, xlen int) {
	GlobalLock()
	scr.Resize(ylen, xlen)
	GlobalUnlock()
}

func MoveScreen(scr *goncurses.Window, ypos int, xpos int) {
	GlobalLock()
	scr.Move(ypos, xpos)
	GlobalUnlock()
}

func BoxScreen(scr *goncurses.Window) {
	GlobalLock()
	err := scr.Border(goncurses.ACS_VLINE, goncurses.ACS_VLINE,
		goncurses.ACS_HLINE, goncurses.ACS_HLINE,
		goncurses.ACS_ULCORNER, goncurses.ACS_URCORNER,
		goncurses.ACS_LLCORNER, goncurses.ACS_LRCORNER)
	if err != nil {
		panic(err)
	}
	GlobalUnlock()
}

func BkgdScreen(scr *goncurses.Window) {
	GlobalLock()
	if color_attr != goncurses.A_NORMAL {
		scr.SetBackground(color_attr | ' ')
	}
	GlobalUnlock()
}
