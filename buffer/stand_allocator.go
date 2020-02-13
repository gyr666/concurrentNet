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
	a 	  Allocator
	load  uint8
}

func (d *divide) Init(a Allocator,load uint8,index uint8){
	d.c = sync.NewCond(&d.l)
	d.load = load
	d.index = index
	d.s= MIN_BLOCK << index
	d.a = a
}

func (d *divide) alloc() ByteBuffer{
	d.Lock()
	d.oper++
	d.give += d.s
	for ;d.first == nil; {
		d.c.Wait()
	}
	v := d.first.(*sByteBuffer)
	d.first = v.addLast(nil)
	// pre allocator
	if d.first == nil {
		go d.backAlloc(true)
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

func (d *divide) backAlloc(release bool) {
	d.Lock()
	result := make([]sByteBuffer,d.load)
	for i := 0;;i++ {
		result[i].Init(d.s,d.a)
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
	sumSize uint64
}

func (s *sByteBuffer) addLast(buffer ByteBuffer) ByteBuffer{
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
	w		sync.WaitGroup
}

func (s *standAllocator) Init(i uint64) error{
	// create alloc index
	s.psize = uint8(i)
	s.min = MIN_BLOCK
	s.max = MIN_BLOCK << s.psize-1
	s.divs = make([]divide,s.psize)
	s.w.Add(int(s.psize))
	for _,v:= range s.divs {
		go func(vl *divide) {
			vl.Init(s, s.psize, uint8(i))
			// at init process, don't need create threads.
			vl.backAlloc(false)
			s.w.Done()
		}(&v)
	}
	s.w.Wait()
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
	release := b.(*sByteBuffer)
	for ;release != nil; {
		go func() { s.divs[release.index].release(release) }()
		release = release.next.(*sByteBuffer)
	}
}

func (s *standAllocator) doAlloc(length uint64) ByteBuffer {
	 r :=position(length)
	 //divide first
	 var b *sByteBuffer = s.divs[r[len(r)-1]].alloc().(*sByteBuffer)
	 l :=b
	 for i := len(r)-2;i>=0;i--{
		 l.next = s.divs[r[i]].alloc()
		 l = l.next.(*sByteBuffer)
	 }
	return b
}

func position(length uint64)[]uint8 {
	return util.IsPow2(length)
}

func (s *standAllocator) AllocSize() uint64 {
	var val uint64 = 0
	for _,v := range s.divs {
		val += v.give
	}
	return val
}

func (b *sByteBuffer) Read(int) ([]byte, error){
	return util.StandRead(0,b.s,b.capital,&b.RP)
}

func (b *sByteBuffer) Write(_b []byte) error {
	//this is too simple
	if util.Int2Uint64(len(_b)) < b.capital-b.WP {
		return util.StandWrite(b.s,b.capital,&b.WP,_b)
	}
	return nil
	//
	//if b.next != nil && util.Int2Uint64(len(_b)) < (b.capital-b.WP)+b.next.Size() {
	//
	//}
}