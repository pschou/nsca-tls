package main

import (
	"bytes"
	"sync"
)

var (
	max_queue_size int64
)

type Buffer struct {
	b bytes.Buffer
	m sync.Mutex
	c chan int
}

func (b *Buffer) ReadString(delim byte) (string, error) {
	if b.b.Len() == 0 || len(b.c) > 0 {
		<-b.c
	}
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadString(delim)
}
func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.b.Len() == 0 {
		<-b.c
	}
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	if b.b.Len() == 0 {
		b.c <- 0
	}
	for b.b.Len()+len(p) > int(max_queue_size) {
		b.b.ReadBytes('\n') // empty out oldest if too big
		//log.Println("Warning: metrics dropped")
	}
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *Buffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}
