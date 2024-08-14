package main

import (
	"bufio"
	"os"
)

type buffer struct {
	f       string
	fd      *os.File
	prev    []string
	blockCh chan int
}

func (this *buffer) init(f string) {
	this.f = f
	fd, err := os.Open(this.f)
	if err == nil {
		this.fd = fd
	}

	this.prev = []string{}
	this.blockCh = make(chan int, 1)
	this.blockCh <- 1
}

func (this *buffer) blockTillReady() {
	<-this.blockCh
}

func (this *buffer) signalBlocked() {
	this.blockCh <- 1
}

func (this *buffer) getMaxLine() int {
	return len(this.prev) - 1
}

func (this *buffer) isDead() bool {
	return this.fd == nil
}

func (this *buffer) update() {
	if this.isDead() {
		return
	}
	this.blockTillReady()
	tmp, _ := this.fd.Seek(0, 1)
	_, _ = this.fd.Seek(0, 0)
	l, err := this.readLines()
	if err == nil {
		this.saveLines(l)
	}
	_, _ = this.fd.Seek(tmp, 0)
	this.signalBlocked()
}

func (this *buffer) readLines() ([]string, error) {
	var ret []string

	scanner := bufio.NewScanner(this.fd)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, scanner.Err()
}

func (this *buffer) saveLines(l []string) {
	this.prev = l
}

func (this *buffer) clear() {
	_, _ = this.fd.Seek(0, 0)
}
