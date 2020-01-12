package test

import (
	"gunplan.top/concurrentNet/buffer"
	"testing"
)

// TestA  Test 必须写，后边的名称随便写
func TestAllocator(t *testing.T) {
	v := buffer.NewBufferAllocator()
	b := v.Alloc(30)
	t.Logf("Size:%d", b.Size())
	b.Release()
	t.Logf("Size:%d", b.Size())

}
