package leakybuf

import "testing"

func TestLeakyBuf(t *testing.T) {
	b := Bytes.Get()
	if b == nil {
		t.Errorf("Got invalid []byte.")
	}
	b = []byte{0x01}
	Bytes.Put(b)
	b = Bytes.Get()
	if int(b[0]) != 0 {
		t.Errorf("It seems that byte used: %d.", len(b))
	}
}

func BenchmarkLeakyBufGetAndPut(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b := Bytes.Get()
			Bytes.Put(b)
		}
	})
}
