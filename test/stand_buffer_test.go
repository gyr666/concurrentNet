package test

import (
	"fmt"
	"testing"

	"gunplan.top/concurrentNet/buffer"
)

// TestA  Test 必须写，后边的名称随便写
func TestStandAllocator(t *testing.T) {
	v := buffer.NewSandBufferAllocator()
	g := v.Alloc(40)
	g.Size()
	for i := 0; i < 40; i++ {
		Equal(v.Alloc(17).Size(), 18, "Size Test")
	}
	b := v.Alloc(41)
	s := "hello eo the world 1234567898 ===UUUU=== u"
	for i := 0;i<10 ; i++ {
		var a =
			[]byte(s)
		b.Write(a)
		bs, _ := b.Read(len(a))
		EqualString(string(bs),s,"RW Test")
		b.Reset()
	}

	b.Release()
	fmt.Printf("operator time : %d\n", v.OperatorTimes())
	fmt.Printf("operator size : %d\n", v.AllocSize())

	t.Logf("Size:%d", v.PoolSize())
}

func TestStandAllocator1(t *testing.T) {
	v := buffer.NewSandBufferAllocator()
	l := v.Alloc(30)
	e := l.Write([]byte{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25})
	l.Write([]byte{26,66})
	c,_ := l.Read(8)
	c,_ = l.Read(8)

	fmt.Print(l,e,c)
}