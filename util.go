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

func GetTerminalInfo() (*winsize, syscall.Errno) {
	ws := &winsize{}
	ret, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(ret) == -1 {
		return nil, errno
	}

	return ws, 0
}

func GetTerminalLines() int {
	GlobalLock()
	ws, errno := GetTerminalInfo()
	GlobalUnlock()

	ret := ws.ws_row
	dbg("LINES", ret, errno)

	return int(ret)
}

func GetTerminalCols() int {
	GlobalLock()
	ws, errno := GetTerminalInfo()
	GlobalUnlock()

	ret := ws.ws_col
	dbg("LINES", ret, errno)

	return int(ret)
}

func Assert(c bool) {
	if !c {
		panic("Assertion")
	}
}
