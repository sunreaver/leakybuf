package leakybuf

import "testing"

func TestLeakyBuffer(t *testing.T) {
	b := Buffers.Get()
	if b == nil {
		t.Errorf("Got invalid buffer.")
	}
	b.Write([]byte{0x01})
	Buffers.Put(b)
	b = Buffers.Get()
	bs := b.Bytes()
	if len(bs) != 0 {
		t.Errorf("bs not reset.")
	}
}

func BenchmarkLeakyBufferGetAndPut(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b := Buffers.Get()
			Buffers.Put(b)
		}
	})
}
