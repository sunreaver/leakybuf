package leakybuf

import (
	"fmt"
	"sync"
)

// Bytes for reuseable []byte. Default size will be 8, maxBuf will be 1024.
var Bytes = NewSyncPool(8, 4096, 2)

// SyncPool is a sync.Pool base slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool
	classesSize []int
	minSize     int
	maxSize     int
}

// NewSyncPool create a sync.Pool base slab allocation memory pool.
// minSize is the smallest chunk size.
// maxSize is the lagest chunk size.
// factor is used to control growth of chunk size.
func NewSyncPool(minSize, maxSize, factor int) *SyncPool {
	n := 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	pool := &SyncPool{
		make([]sync.Pool, n),
		make([]int, n),
		minSize, maxSize,
	}
	n = 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		pool.classesSize[n] = chunkSize
		pool.classes[n].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]byte, size)
				return &buf
			}
		}(chunkSize)
		n++
	}
	return pool
}

// Get []byte
// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (pool *SyncPool) Get(sizes ...int) []byte {
	var size int
	if len(sizes) == 0 {
		size = pool.minSize
	} else {
		size = sizes[0]
	}
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]byte)
				fmt.Println(size, cap(*mem))
				return (*mem)[:size]
			}
		}
	}
	return make([]byte, size)
}

// Put []byte
// Free release a []byte that alloc from Pool.Alloc.
func (pool *SyncPool) Put(mem []byte) {
	if size := cap(mem); size <= pool.maxSize && size >= pool.minSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				pool.classes[i].Put(&mem)
				return
			}
		}
	}
}
