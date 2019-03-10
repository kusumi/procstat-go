// +build stdout

package main

import (
	"fmt"
	"time"

	"github.com/rthornton128/goncurses"
)

func StringToColor(arg string) int16 {
	return -1
}

func InitScreen(fg, bg int16) {
}

func CleanupScreen() {
}

func ReadIncoming() int {
	time.Sleep(200 * time.Millisecond)
	return 0
}

func ClearTerminal() {
}

func FlashTerminal() {
}

func AllocScreen(ylen, xlen, ypos, xpos int) *goncurses.Window {
	return nil
}

func DeleteScreen(scr *goncurses.Window) {
}

func PrintScreen(scr *goncurses.Window, y int, x int, standout bool, s string) {
	fmt.Println(s)
}

func RefreshScreen(scr *goncurses.Window) {
}

func EraseScreen(scr *goncurses.Window) {
}

func ResizeScreen(scr *goncurses.Window, ylen int, xlen int) {
}

func MoveScreen(scr *goncurses.Window, ypos int, xpos int) {
}

func BoxScreen(scr *goncurses.Window) {
}

func BkgdScreen(scr *goncurses.Window) {
}
