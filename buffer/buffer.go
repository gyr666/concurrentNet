package buffer

import (
	"container/list"
	"gunplan.top/concurrentNet/util"
	"sync"
)

type OperatorMode int

const (
	READ  OperatorMode = 1
	WRITE OperatorMode = 1 << 1
)

type Allocator interface {
	Develop
	Alloc(uint64) ByteBuffer
	release(ByteBuffer)
	UsedSize() int
}

type ByteBuffer interface {
	Develop
	Release()
	Size() uint64
}

func NewBufferAllocator() Allocator {
	a := allocatorImpl{}
	if err := a.Init(); err == nil {
		return &a
	}
	return nil

}

type allocatorImpl struct {
	unUsed *list.List
	used   *list.List
	l      sync.Mutex
}

func (a *allocatorImpl) Init() error {
	a.unUsed = list.New()
	a.used = list.New()
	return nil
}

func (a *allocatorImpl) Destroy() error {
	a.unUsed = nil
	a.used = nil
	//help gc
	return nil
}
func (a *allocatorImpl) UsedSize() int {
	return a.used.Len()
}

func (a *allocatorImpl) Alloc(length uint64) ByteBuffer {
	a.l.Lock()
	if bf, err := a.findByUnusedList(length); err == nil && bf != nil {
		return bf
	}
	b := byteBufferImpl{a: a, spaceLength: length}
	b.Init()
	a.addToUsedList(&b)
	a.l.Unlock()
	return &b
}

func (a *allocatorImpl) removeUsed(buffer ByteBuffer) {
	util.Traverse(a.used, func(element *list.Element) bool {
		return (element.Value.(ByteBuffer)) == buffer
	})
}

func (a *allocatorImpl) findByUnusedList(i uint64) (ByteBuffer, error) {
	var buffer ByteBuffer = nil
	util.Traverse(a.unUsed, func(element *list.Element) bool {
		if ((element.Value.(ByteBuffer)).Size() - i) < 10 {
			buffer = element.Value.(ByteBuffer)
			buffer.Destroy()
			return true
		}
		return false
	})
	return buffer, nil
}
func (a *allocatorImpl) addToUsedList(buffer ByteBuffer) {
	a.used.PushBack(buffer)
}

func (a *allocatorImpl) release(buffer ByteBuffer) {
	a.l.Lock()
	a.unUsed.PushBack(buffer)
	a.removeUsed(buffer)
	a.l.Unlock()
}

type byteBufferImpl struct {
	a           Allocator
	spaceLength uint64
	RP          uint64
	WP          uint64
	s           []byte
	mode        OperatorMode
}

func (a *byteBufferImpl) Init() error {
	a.s = make([]byte, a.spaceLength)
	return nil
}

func (a *byteBufferImpl) Destroy() error {
	a.RP = 0
	a.WP = 0
	a.spaceLength = 0
	a.s = nil
	// help gc
	return nil
}
func (b *byteBufferImpl) Size() uint64 {
	return b.spaceLength
}

func (b *byteBufferImpl) Release() {
	b.a.release(b)
}

func (b *byteBufferImpl) Mode() OperatorMode {
	return b.mode
}
