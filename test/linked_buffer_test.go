package test

import (
	"fmt"
	"testing"

	"gunplan.top/concurrentNet/buffer"
)

func TestAllocator(t *testing.T) {
	v := buffer.NewLikedBufferAllocator()
	b := v.Alloc(30)
	t.Logf("Size:%d", b.Size())
	b.Release()
	b = v.Alloc(30)
	t.Logf("Size:%d", b.Size())
	b.Release()
	b = v.Alloc(30)
	t.Logf("Size:%d", b.Size())
	b.Release()
	b = v.Alloc(30)
	t.Logf("Size:%d", b.Size())
	var a = []byte("hello")
	b.Write(a)
	bs, err := b.Read(len(a))
	if err == nil {
		fmt.Println(string(bs))
	}
	b.Release()

	t.Logf("Size:%d", v.PoolSize())

}
