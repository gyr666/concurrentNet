package test

import (
	"fmt"
	"gunplan.top/concurrentNet/buffer"
	"testing"
)

// TestA  Test 必须写，后边的名称随便写
func TestStandAllocator(t *testing.T) {
	v := buffer.NewSandBufferAllocator()
	b := v.Alloc(32)
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
