package buffer

import (
	"gunplan.top/concurrentNet/util"
	"sync"
)

func NewLikedBufferAllocator() Allocator {
	a := allocatorImpl{}
	if err := a.Init(0); err == nil {
		return &a
	}
	return nil
}

type allocatorImpl struct {
	counter util.DynamicCounter
	unUsed  *util.Skiplist
	l       sync.Mutex
}

func (a *allocatorImpl) Init(uint64) error {
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
	}
	bf := new(byteBufferImpl)
	bf.Init(length, a)
	return bf
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
	a.counter.Push(-buffer.Size())
	a.unUsed.Insert(buffer.Size(), buffer)
	a.l.Unlock()
	go a.dynamicShrink()
}

func (a *allocatorImpl) dynamicShrink() {

}

type byteBufferImpl struct {
	BaseByteBuffer
}

func (a *allocatorImpl) OperatorTimes() uint64 {
	return a.counter.Size()
}

func (a *allocatorImpl) AllocSize() uint64 {
	return a.counter.Use()
}