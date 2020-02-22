package util

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type Spinlock uint32

func NewSpinlock() sync.Locker {
	return new(Spinlock)
}

func (s *Spinlock) Lock() {
	if !atomic.CompareAndSwapUint32((*uint32)(s), 0, 1) {
		runtime.Gosched()
	}
}

func (s *Spinlock) Unlock() {
	atomic.StoreUint32((*uint32)(s), 0)
}
