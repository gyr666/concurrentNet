package buffer

import (
	"gunplan.top/concurrentNet/util"
	"sync"
)

func NewBufferAllocator() Allocator {
	a := allocatorImpl{}
	if err := a.Init(); err == nil {
		return &a
	}
	return nil
}

type allocatorImpl struct {
	counter util.DynamicCounter
	unUsed  *util.Skiplist
	l       sync.Mutex
}

func (a *allocatorImpl) Init() error {
	a.unUsed = util.NewSkipList()
	a.counter = util.NewCounter()
	a.counter.Boot()
	return nil
}

func (a *allocatorImpl) PoolSize() uint64 {
	return a.counter.Size()
}

func (a *allocatorImpl) Destroy() error {
	a.unUsed = nil
	a.counter = nil
	//help gc
	return nil
}

func (a *allocatorImpl) Alloc(length uint64) ByteBuffer {
	a.counter.Push(length)
	a.l.Lock()
	defer a.l.Unlock()
	if bf := a.findByUnusedList(length); bf != nil {
		return bf
	} else {
		bf := byteBufferImpl{BaseByteBuffer: BaseByteBuffer{capital: length}, a: a}
		bf.Init()
		return &bf
	}
}	

func (a *allocatorImpl) findByUnusedList(i uint64) ByteBuffer {
	k, v := a.unUsed.Search(i)
	if k != 0 {
		a.unUsed.Delete(k)
		return v.(ByteBuffer)
	}
	return nil
}

func (a *allocatorImpl) release(buffer ByteBuffer) {
	a.l.Lock()
	a.unUsed.Insert(buffer.Size(), buffer)
	a.l.Unlock()
	go a.dynamicShrink()
}

func (a *allocatorImpl) dynamicShrink() {

}

type byteBufferImpl struct {
	BaseByteBuffer
	a    Allocator
	s    []byte
	mode OperatorMode
}

func (a *byteBufferImpl) Init() error {
	a.s = make([]byte, a.capital)
	return nil
}

func (a *byteBufferImpl) Destroy() error {
	a.RP = 0
	a.WP = 0
	return nil
}

func (b *byteBufferImpl) Release() {
	b.a.release(b)
}


