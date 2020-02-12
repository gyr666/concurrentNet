// A fast memory allocator support high concurrency
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
	first ByteBuffer
	l     sync.Mutex
	s     uint64
	c     *sync.Cond
	index uint8
	give  uint64
	oper  uint64
	load  uint8
}

func (d *divide) Init(load uint8,index uint8){
	d.c = sync.NewCond(&d.l)
	d.load = load
	d.index = index
	d.s= MIN_BLOCK << index
}

func (d *divide) alloc(s Allocator,next bool) ByteBuffer{
	d.Lock()
	d.oper++
	d.give += d.s
	for ;d.first == nil; {
		d.c.Wait()
	}
	v := d.first.(*sByteBuffer)
	var nb ByteBuffer
	if next {
		nb = s.Alloc(d.s>>1)
	}
	d.first = v.c434bb(nb)
	// pre allocator
	if d.first == nil {
		go d.backAlloc(s,true)
	}
	d.Unlock()
	return v
}

func (d *divide) release(buffer *sByteBuffer){
	d.l.Lock()
	defer d.l.Unlock()
	d.oper++
	d.give -= d.s
	buffer.next = d.first
	d.first = buffer
}

func (d *divide) Lock(){
	d.l.Lock()
}

func (d *divide) Unlock(){
	d.l.Unlock()
}

func (d *divide) setLink(sta *sByteBuffer){
	d.first = sta
}

func (d *divide) backAlloc(s Allocator,release bool) {
	d.Lock()
	result := make([]sByteBuffer,d.load)
	for i := 0;;i++ {
		result[i].Init(d.s,s)
		result[i].index = d.index
		if i == len(result)-1 {
			break
		}
		result[i].next = &result[i+1]
	}
	d.setLink(&result[0])
	d.Unlock()
	if release {
		d.c.Broadcast()
	}
}

type sByteBuffer struct {
	BaseByteBuffer
	next	ByteBuffer
	index	uint8
}

func (s *sByteBuffer) c434bb(buffer ByteBuffer) ByteBuffer{
	p := s.next
	s.next = buffer
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
	s.min = MIN_BLOCK
	s.max = MIN_BLOCK << s.psize-1
	s.divs = make([]divide,s.psize)
	for i,_:= range s.divs {
		s.divs[i].Init(s.psize,uint8(i))
		// at init process, don't need create threads.
		s.divs[i].backAlloc(s,false)
	}
	return nil
}

func (s *standAllocator) OperatorTimes() uint64 {
	var val uint64 = 0
	for _,v := range s.divs {
		val += v.oper
	}
	return val
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
	s.divs[b.(*sByteBuffer).index].release(b.(*sByteBuffer))
}

func (s *standAllocator) doAlloc(length uint64) ByteBuffer {
	index,next:= position(length)
	return s.divs[index].alloc(s,next)
}

func position(length uint64) (uint8,bool) {
	v,ok := util.IsPow2(length)
	if !ok && length& (1<<(v-1)) != 0{
		return uint8(v),false
	} else if !ok && length& (1<<(v-1)) == 0 {
		return uint8(v) - 1 ,true
	}
	return uint8(v) - 1,false
}

func (s *standAllocator) AllocSize() uint64 {
	var val uint64 = 0
	for _,v := range s.divs {
		val += v.give
	}
	return val
}