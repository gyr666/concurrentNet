package test

import (
	"fmt"
	"gunplan.top/concurrentNet/buffer"
	"testing"
)

// TestA  Test 必须写，后边的名称随便写
func TestStandAllocator(t *testing.T) {
	v := buffer.NewSandBufferAllocator()
	for i:=0;i<40;i++{
		Equal(v.Alloc(30).Size(),32,"Size Test")
	}
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
	fmt.Printf("operator time : %d\n",v.OperatorTimes())
	fmt.Printf("operator size : %d\n",v.AllocSize())

	t.Logf("Size:%d", v.PoolSize())
}