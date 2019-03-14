package main

import (
	"bufio"
	"os"
)

type Buffer struct {
	f    string
	fd   *os.File
	prev []string
	l_ch chan bool
}

func (this *Buffer) Init(f string) {
	this.f = f
	fd, err := os.Open(this.f)
	if err == nil {
		this.fd = fd
	}

	this.prev = []string{}
	this.l_ch = make(chan bool, 1)
	this.l_ch <- true
}

func (this *Buffer) BlockTillReady() {
	<-this.l_ch
}

func (this *Buffer) SignalBlocked() {
	this.l_ch <- true
}

func (this *Buffer) GetMaxLine() int {
	return len(this.prev) - 1
}

func (this *Buffer) IsDead() bool {
	return this.fd == nil
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
	this.fd.Seek(0, 0)
}