package main

import (
	"bufio"
	"os"
)

type Buffer struct {
	f    string
	fd   *os.File
	prev []string
	l_ch chan int
}

func (this *Buffer) Init(f string) {
	this.f = f
	fd, err := os.Open(this.f)
	if err == nil {
		this.fd = fd
	}

	this.prev = []string{}
	this.l_ch = make(chan int, 1)
	this.l_ch <- 1
}

func (this *Buffer) BlockTillReady() {
	<-this.l_ch
}

func (this *Buffer) SignalBlocked() {
	this.l_ch <- 1
}

func (this *Buffer) GetMaxLine() int {
	return len(this.prev) - 1
}

func (this *Buffer) IsDead() bool {
	return this.fd == nil
}

func (this *Buffer) Update() {
	if this.IsDead() {
		return
	}
	this.BlockTillReady()
	tmp, _ := this.fd.Seek(0, 1)
	_, _ = this.fd.Seek(0, 0)
	l, err := this.ReadLines()
	if err == nil {
		this.SaveLines(l)
	}
	_, _ = this.fd.Seek(tmp, 0)
	this.SignalBlocked()
}

func (this *Buffer) ReadLines() ([]string, error) {
	var ret []string

	scanner := bufio.NewScanner(this.fd)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, scanner.Err()
}

func (this *Buffer) SaveLines(l []string) {
	this.prev = l
}

func (this *Buffer) Clear() {
	_, _ = this.fd.Seek(0, 0)
}
