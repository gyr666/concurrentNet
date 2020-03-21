// A fast memory allocator support high concurrency
package buffer

import (
	"errors"
	"sync"
	"time"

	"gunplan.top/concurrentNet/util"
)

const MinBlock = 2

func NewSandBufferAllocator() Allocator {
	alloc := standAllocator{}
	_ = alloc.Init(20)
	return &alloc
}

func (s *sByteBuffer) Release() {
	s.a.release(s)
}

type divide struct {
	first ByteBuffer
	l     sync.Mutex
	s     uint64
	c     *sync.Cond
	index uint8
	give  uint64
	oper  uint64
	a     Allocator
	// load is a dynamic value
	load uint8
}

func (d *divide) Init(a Allocator, load uint8, index uint8) {
	d.c = sync.NewCond(&d.l)
	d.load = load
	d.index = index
	d.s = MinBlock << index
	d.a = a
}

func (d *divide) Destroy() {
	d.first = nil
	// help gc
}

func (d *divide) alloc() ByteBuffer {
	d.Lock()
	d.oper++
	d.give += d.s
	for d.first == nil {
		d.c.Wait()
	}
	v := d.first.(*sByteBuffer)
	d.first = v.addLast(nil)
	// pre allocator
	if d.first == nil {
		go d.backAlloc()
	}
	d.Unlock()
	return v
}

func (d *divide) release(buffer *sByteBuffer) {
	d.l.Lock()
	defer d.l.Unlock()
	d.oper++
	d.give -= d.s
	buffer.next = d.first
	d.first = buffer
}

func (d *divide) Lock() {
	d.l.Lock()
}

func (d *divide) Unlock() {
	d.l.Unlock()
}

func (d *divide) setLink(sta *sByteBuffer) {
	d.first = sta
}

func (d *divide) backAlloc() {
	d.Lock()
	result := make([]sByteBuffer, d.load)
	for i := 0; i < len(result); i++ {
		result[i].init(d.s, d.a, d.index)
		if i == len(result)-1 {
			result[i].next = nil
		} else {
			result[i].next = &result[i+1]
		}
	}
	d.setLink(&result[0])
	d.Unlock()
	d.c.Broadcast()
}

type sByteBuffer struct {
	BaseByteBuffer
	next    ByteBuffer
	index   uint8
	sumSize uint64
	active  []bool
}

func (s *sByteBuffer) init(size uint64, all Allocator, index uint8) {
	s.BaseByteBuffer.Init(size, all)
	s.sumSize = size
	s.active = make([]bool, 2)
	s.index = index
}

func (s *sByteBuffer) Size() uint64 {
	return s.sumSize
}

func (s *sByteBuffer) addLast(buffer ByteBuffer) ByteBuffer {
	p := s.next
	s.next = buffer
	return p
}

type standAllocator struct {
	divs  []divide
	min   uint64
	psize uint8
	load  uint8
	r     bool
	max   uint64
	regs  int64
	w     sync.WaitGroup
}

func (s *standAllocator) Init(i uint64) error {
	// create alloc index
	s.r = true
	s.psize = uint8(i)
	s.min = MinBlock
	s.max = MinBlock<<s.psize - 1
	s.divs = make([]divide, s.psize)
	s.w.Add(int(s.psize))
	for i := range s.divs {
		go func(vl *divide, index int) {
			vl.Init(s, s.psize, uint8(index))
			// at init process, don't need create threads.
			vl.backAlloc()
			s.w.Done()
		}(&s.divs[i], i)
	}
	go s.dynamicRegulate()
	s.w.Wait()
	return nil
}

func (s *standAllocator) OperatorTimes() uint64 {
	var val uint64 = 0
	for _, v := range s.divs {
		val += v.oper
	}
	return val
}

func (s *standAllocator) Destroy() error {
	s.r = false
	for i := range s.divs {
		s.divs[i].Destroy()
	}
	return nil
}

func (s *standAllocator) Alloc(length uint64) ByteBuffer {
	if length > s.max || length == 0 {
		return nil
	}
	return s.doAlloc(length)
}

func (s *standAllocator) PoolSize() uint64 {
	return util.Int2Uint64(int(s.psize))
}

func (s *standAllocator) release(b ByteBuffer) {
	release := b
	next := release
	for next != nil {
		release = release.(*sByteBuffer).next
		go func(n *sByteBuffer) { s.divs[n.index].release(n) }(next.(*sByteBuffer))
		next = release
	}
}

func (s *standAllocator) doAlloc(length uint64) ByteBuffer {
	r := position(length)
	//divide first
	var b = s.divs[r[len(r)-1]].alloc().(*sByteBuffer)
	l := b
	for i := len(r) - 2; i >= 0; i-- {
		l.next = s.divs[r[i]].alloc()
		l = l.next.(*sByteBuffer)
		b.sumSize += l.capital
	}
	return b
}

func position(length uint64) []uint8 {
	return util.IsPow2(length)
}

func (s *standAllocator) AllocSize() uint64 {
	var val uint64 = 0
	for _, v := range s.divs {
		val += v.give
	}
	return val
}

func (s *standAllocator) dynamicRegulate() {
	c := util.NewCounter()
	c.Boot()
	for s.r {
		for i := 0; i < len(s.divs); i++ {
			c.Push(s.divs[i].oper)
		}
		time.Sleep(time.Second * 2)
		ave := c.Ave()
		for i := 0; i < len(s.divs); i++ {
			if s.divs[i].oper <= ave {
				s.divs[i].load--
			} else {
				s.divs[i].load++
			}
		}
		c.Reset()
	}
}

func (s *sByteBuffer) Read(len uint64) ([]byte, error) {
	if s.sumSize-s.RP < len {
		return nil, errors.New(util.INDEX_OUTOF_BOUND)
	}
	send := make([]byte, len)
	s.Read0(len, 0, send)
	return send, nil
}

func (s *sByteBuffer) ReadAll() ([]byte, error) {
	return nil, nil
}

func (s *sByteBuffer) Read0(len, pos uint64, send []byte) {
	i := pos
	for ; i < len; i++ {
		if s.RP == s.capital {
			break
		}
		send[i] = util.ReadOne(s.s, &s.RP)
	}
	if i < len {
		s.next.(*sByteBuffer).Read0(len, i, send)
	}
}

func (s *sByteBuffer) Write(_b []byte) error {
	if s.sumSize-s.WP < uint64(len(_b)) {
		return errors.New(util.INDEX_OUTOF_BOUND)
	}
	s.Write0(_b, 0)
	return nil
}

func (s *sByteBuffer) Write0(_b []byte, position uint64) {
	i := position
	for ; i < uint64(len(_b)); i++ {
		if s.WP == s.capital {
			break
		}
		util.WriteOne(s.s, _b[i], &s.WP)
	}
	if i != uint64(len(_b)) {
		s.next.(*sByteBuffer).Write0(_b, i)
	}
}

func (s *sByteBuffer) AvailableReadSum() uint64 {
	return s.globalWP(0) - s.globalRP(0) - 1
}

func (s *sByteBuffer) globalWP(now uint64) uint64 {
	if s.WP != s.capital || s.next == nil {
		return now + s.WP
	} else {
		return s.next.(*sByteBuffer).globalWP(now + s.capital)
	}
}

func (s *sByteBuffer) globalRP(now uint64) uint64 {
	if s.WP != s.capital || s.next == nil {
		return now + s.RP
	} else {
		return s.next.(*sByteBuffer).globalRP(now + s.capital)
	}
}

func (s *sByteBuffer) FastMoveOut() *[]byte {
	panic("sByteBuffer method `FastMoveOut` not support!")
	return nil
}

func (s *sByteBuffer) FastMoveIn(*[]byte) {
	panic("sByteBuffer method `FastMoveIn` not support!")
}
