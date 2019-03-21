// +build stdout

package main

import (
	"fmt"
	"time"
)

const (
	KEY_ERR    = -1
	KEY_UP     = 1
	KEY_DOWN   = 2
	KEY_LEFT   = 3
	KEY_RIGHT  = 4
	KEY_RESIZE = 5
)

type Screen struct {
}

func StringToColor(arg string) int16 {
	return -1
}

func InitScreen(fg, bg int16) error {
	return nil
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

func AllocScreen(ylen, xlen, ypos, xpos int) *Screen {
	return nil
}

func (this *Screen) Delete() {
}

func (this *Screen) Print(y, x int, standout bool, s string) {
	fmt.Println(s)
}

func (this *Screen) Refresh() {
}

func (this *Screen) Erase() {
}

func (this *Screen) Resize(ylen, xlen int) {
}

func (this *Screen) Move(ypos, xpos int) {
}

func (this *Screen) Box() {
}

func (this *Screen) Bkgd() {
}
