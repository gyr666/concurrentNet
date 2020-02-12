package buffer

import (
	"gunplan.top/concurrentNet/util"
	"sync"
)
const MIN_BLOCK = 2

func NewSandBufferAllocator() Allocator{
	alloc := standAllocator{}
	alloc.Init(20)
	return &alloc
}

func (b *sByteBuffer) Release() {
	b.a.release(b)
}

type divide struct {
	first	ByteBuffer
	last	ByteBuffer
	l		sync.Mutex
	s		uint64
	c       *sync.Cond
}

func (d *divide) Init(s uint64){
	d.c = sync.NewCond(&d.l)
	d.s=s
}

func (d *divide) Lock(){
	d.l.Lock()
}

func (d *divide) Unlock(){
	d.l.Unlock()
}

func (d *divide) setLink(sta,end *sByteBuffer){
	d.first = sta
	d.last = end
}

type sByteBuffer struct {
	BaseByteBuffer
	next	ByteBuffer
	index	uint8
}

func (s *sByteBuffer) c434bb() ByteBuffer{
	p := s.next
	s.next = nil
	return p
}

type standAllocator struct {
	divs	[]divide
	min		uint64
	psize	uint8
	load	uint8
	max		uint64
}

func (s *standAllocator) Init(i uint64) error{
	// create alloc index
	s.psize = uint8(i)
	s.load = s.psize
	s.min = MIN_BLOCK
	s.max = MIN_BLOCK << s.psize-1;
	s.divs = make([]divide,s.psize)
	for i,_:= range s.divs {
		s.divs[i].Init(2<<i)
		// at init process, don't need create threads.
		s.backAlloc(util.Int2UInt8(i),false)
	}
	return nil
}

func (s *standAllocator) Destroy() error{
	return nil
}

func (s *standAllocator) Alloc(length uint64) ByteBuffer {
	if(length > s.max) {
		return nil
	}
	return s.doAlloc(length)
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
	index := position(length)
	block := &s.divs[index]
	block.Lock()
	for ;block.first == nil; {
		block.c.Wait()
	}
	v := block.first.(*sByteBuffer)
	block.first = v.c434bb()
	// pre allocator
	if block.first == block.last {
		go s.backAlloc(index,true)
	}
	block.Unlock()
	return v
}

func position(length uint64) uint8 {
	v,ok := util.IsPow2(length)
	if !ok {
		return uint8(v)
	}
	return uint8(v) - 1
}

//this process is very important
func (s *standAllocator) backAlloc(index uint8,release bool) {
	s.divs[index].Lock()
	defer s.divs[index].Unlock()
	result := make([]sByteBuffer,s.load)
	for i := 0;i<len(result)-1 ;i++ {
		result[i].Init(s.divs[index].s,s)
		result[i].index = index
		result[i].next = &result[i+1]
	}
	s.divs[index].setLink(&result[0],&result[s.load-1])
	if release {
		s.divs[index].c.Broadcast()
	}
}
