package main

import (
	"syscall"
	"time"
	"unsafe"
)

var globalCh = make(chan int, 1)

func initLock() {
	globalCh <- 1
}

func cleanupLock() {
	close(globalCh)
}

func globalLock() {
	<-globalCh
}

func globalUnlock() {
	globalCh <- 1
}

type winsize struct {
	wsRow    uint16
	wsCol    uint16
	wsXpixel uint16
	wsYpixel uint16
}

func getTerminalInfo() (*winsize, syscall.Errno) {
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

func getTerminalLines() int {
	globalLock()
	ws, errno := getTerminalInfo()
	globalUnlock()

	ret := ws.wsRow
	dbg("LINES", ret, errno)
	if errno != 0 {
		panic(errno)
	}

	return int(ret)
}

func getTerminalCols() int {
	globalLock()
	ws, errno := getTerminalInfo()
	globalUnlock()

	ret := ws.wsCol
	dbg("COLS", ret, errno)
	if errno != 0 {
		panic(errno)
	}

	return int(ret)
}

func getSecond(t int) time.Duration {
	return time.Duration(t) * time.Second
}

func getMillisecond(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

func assert(c bool) {
	kassert(c, "Assert failed")
}

func kassert(c bool, err interface{}) {
	if !c {
		panic(err)
	}
}
