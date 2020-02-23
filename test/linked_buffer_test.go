package test

import (
	"gunplan.top/concurrentNet/buffer"
	"math/rand"
	"testing"
)

var all = buffer.NewLikedBufferAllocator

func TestAllocatorCreate(t *testing.T) {
	a := all()
	a.Destroy()
}

func TestAllocatorAllocate(t *testing.T) {
	TestAllocatorCreate(t)
	a := all()
	for i := 0; i < 200; i++ {
		v := rand.Int31n(10000)
		b := a.Alloc(uint64(v))
		Equal(b.Size(), uint64(v), "size test")
		b.Release()
	}
	t.Log(a.OperatorTimes())
	a.Destroy()
}

func TestBufferRw(t *testing.T) {
	a := all()
	b := a.Alloc(32)
	s := "hello world"
	b.Write([]byte(s))
	d, _ := b.Read(uint64(len(s)))
	EqualString(string(d), s, "rw test")
	a.Destroy()
}
