package main

import (
	"syscall"
	"time"
	"unsafe"
)

var gl_ch = make(chan int, 1)

func InitLock() {
	gl_ch <- 1
}

func CleanupLock() {
	close(gl_ch)
}

func GlobalLock() {
	<-gl_ch
}

func GlobalUnlock() {
	gl_ch <- 1
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

func GetSecond(t int) time.Duration {
	return time.Duration(t) * time.Second
}

func GetMillisecond(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

func Assert(c bool) {
	if !c {
		panic("Assertion")
	}
}
