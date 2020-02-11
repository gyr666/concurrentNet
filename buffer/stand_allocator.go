package buffer

import (
	"gunplan.top/concurrentNet/util"
	"sync"
)

func NewSandBufferAllocator() Allocator{
	return &standAllocator{}
}
type divide struct {
	first	ByteBuffer
	last	ByteBuffer
	s		uint64
}

type sByteBuffer struct {
	BaseByteBuffer
	next	ByteBuffer
	index	uint8
}

func (c *sByteBuffer) Init() error{
	return nil
}

func (c *sByteBuffer) Destroy() error{
	return nil
}


func (c *sByteBuffer) Release() {

}



type standAllocator struct {
	divs	[]divide
	min		uint64
	psize	uint8
	load	uint8
	max		uint64
	l	sync.Mutex
}

func (s *standAllocator) Init() error{
	// create alloc index
	s.divs = make([]divide,s.psize)
	for i,_:= range s.divs {
		s.divs[i].s = 2 << i
		// at init process, don't need create threads.
		s.backAlloc(util.Int2UInt8(i))
	}
	s.min = 2;
	s.max = 2 << p.size-1;
	return nil
}

func (s *standAllocator) Destroy() error{
	return nil
}

func (s *standAllocator) Alloc(length uint64) ByteBuffer {
	if(length > s.max) {
		return nil
	}
	return s.doAlloc(length);
}

func (s *standAllocator) PoolSize() uint64 {
	return util.Int2Uint64(int(s.psize))
}

func (s *standAllocator) release(b ByteBuffer) {
	bb ,ok:= b.(*sByteBuffer)
	if !ok {
		panic("release error!")
	}
	bb.next = s.divs[bb.index].first
	s.divs[bb.index].first = bb
}

func (s *standAllocator) doAlloc(length uint64) ByteBuffer {
	index := postition(length)
	if s.divs[index].first == s.divs[index].last {
		go s.backAlloc(index)
	}
	s.l.Lock()
	v := s.divs[index].first.(*sByteBuffer)
	s.divs[index].first = v.next
	s.l.Unlock();
	return v
}

func postition(length uint64) uint8 {
	 v,ok := util.IsPow2(length)
	if ok {
		return v
	}
	return v+1
}

//this process is very important
func (s *standAllocator) backAlloc(index uint8) ByteBuffer {
	s.l.Lock()
	defer s.l.Unlock()
	result := make([]sByteBuffer,s.load)
	for i,_ := range result {
		result[i].next = &result[i+1]
	}
	s.divs[index].first = result
	s.divs[index].last  = result + s.load
}


