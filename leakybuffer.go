package leakybuf

import (
	"bytes"
	"sync"
)

// Buffers for reuseable bytes.Buffer.
var Buffers = newLeakyBuffer()

// LeakyBuffer for bytes.Buffer
type LeakyBuffer struct {
	// freeList chan *bytes.Buffer
	bufPool sync.Pool
}

func newLeakyBuffer() *LeakyBuffer {
	return &LeakyBuffer{
		// make(chan *bytes.Buffer, maxSize),
		bufPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get returns a buffer from the leaky buffer or create a new buffer.
func (lb *LeakyBuffer) Get() (b *bytes.Buffer) {
	// select {
	// case b = <-lb.freeList:
	// 	b.Reset()
	// default:
	// 	b = bytes.NewBuffer(make([]byte, 0, maxBufferSize))
	// }
	b = lb.bufPool.Get().(*bytes.Buffer)
	return
}

// Put add the buffer into the free buffer pool for reuse. Panic if the buffer
// size is not the same with the leaky buffer's. This is intended to expose
// error usage of leaky buffer.
func (lb *LeakyBuffer) Put(b *bytes.Buffer) {
	// select {
	// case lb.freeList <- b:
	// default:
	// }
	b.Reset()
	lb.bufPool.Put(b)
	return
}
