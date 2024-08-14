//go:build stdout

package main

import (
	"fmt"
	"time"
)

type screen struct {
}

const (
	KEY_ERR    = 0xDEAD
	KEY_UP     = KEY_ERR + 1
	KEY_DOWN   = KEY_ERR + 2
	KEY_LEFT   = KEY_ERR + 3
	KEY_RIGHT  = KEY_ERR + 4
	KEY_RESIZE = KEY_ERR + 5
)

func keyCtrl(x int) int {
	return x & 0x1F
}

func stringToColor(arg string) int16 {
	return -1
}

func initScreen(fg, bg int16) error {
	return nil
}

func cleanupScreen() {
}

func readIncoming() int {
	time.Sleep(200 * time.Millisecond)
	return 0
}

func clearTerminal() {
}

func flashTerminal() {
}

func allocScreen(ylen, xlen, ypos, xpos int) screen {
	return screen{}
}

func (this *screen) delete() {
}

func (this *screen) print(y, x int, standout bool, s string) {
	fmt.Println(s)
}

func (this *screen) refresh() {
}

func (this *screen) erase() {
}

func (this *screen) resize(ylen, xlen int) {
}

func (this *screen) move(ypos, xpos int) {
}

func (this *screen) box() {
}

func (this *screen) bkgd() {
}
