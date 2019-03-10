package main

import (
	"syscall"
	"unsafe"
)

var gl_ch = make(chan bool, 1)

func InitLock() {
	gl_ch <- true
}

func CleanupLock() {
	close(gl_ch)
}

func GlobalLock() {
	<-gl_ch
}

func GlobalUnlock() {
	gl_ch <- true
}

type winsize struct {
	ws_row    uint16
	ws_col    uint16
	ws_xpixel uint16
	ws_ypixel uint16
}

func GetTerminalInfo() *winsize {
	ws := &winsize{}
	ret, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(ret) == -1 {
		panic(errno)
	}

	return ws
}

func GetTerminalLines() int {
	GlobalLock()
	ws := GetTerminalInfo()
	GlobalUnlock()

	ret := ws.ws_row
	dbg("LINES", ret)

	return int(ret)
}

func GetTerminalCols() int {
	GlobalLock()
	ws := GetTerminalInfo()
	GlobalUnlock()

	ret := ws.ws_col
	dbg("COLS", ret)

	return int(ret)
}

func Assert(c bool) {
	if !c {
		panic("Assertion")
	}
}
