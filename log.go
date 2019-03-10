package main

import (
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

var (
	linit bool = false
	lfd   *os.File
)

func InitLog() {
	if !opt.debug {
		return
	}

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	f := path.Join(u.HomeDir, ".procstat.log")
	lfd, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	GlobalLock()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(lfd)
	GlobalUnlock()

	linit = true
	dbg(strings.Repeat("=", 20))
	dbg(f, lfd)
}

func CleanupLog() {
	if !opt.debug {
		return
	}

	GlobalLock()
	lfd.Close()
	GlobalUnlock()

	linit = false
}

func dbg(args ...interface{}) {
	if !opt.debug {
		return
	}

	Assert(linit)
	GlobalLock()
	log.Println(args...)
	GlobalUnlock()
}

func dbgf(f string, args ...interface{}) {
	if !opt.debug {
		return
	}

	Assert(linit)
	GlobalLock()
	log.Printf(f, args...)
	GlobalUnlock()
}
