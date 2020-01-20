package disruptor

import (
	"sync/atomic"
)

type disruptorImpl struct {
	buffer []interface{}
	ps     Sequence
	cs     Sequence
}

type sequenceImpl struct {
	max uint64
	now uint64
}

func (s *sequenceImpl) Next() (uint64, bool) {
	if s.now == s.max {
		s.now = 0
	} else {
		s.now++
	}
	return s.now, true
}

type sequenceCASImpl struct {
	max uint64
	now uint64
}

func (s *sequenceCASImpl) Next() (uint64, bool) {
	if s.now == s.max && atomic.CompareAndSwapUint64(&s.now, s.now, 0) {
		return s.now, true
	} else if atomic.CompareAndSwapUint64(&s.now, s.now, s.now+1) {
		return s.now, true
	} else {
		return 0, false
	}

}

func (d *disruptorImpl) Init(n uint64) {
	d.ps = &sequenceImpl{max: n}
	d.cs = &sequenceCASImpl{max: n}
	d.buffer = make([]interface{}, n)

}

func (d *disruptorImpl) Push(i interface{}) {
	n, _ := d.ps.Next()
	d.buffer[n] = i
}

func (d *disruptorImpl) Poll() interface{} {
	if v, t := d.cs.Next(); t {
		return d.buffer[v]
	} else {
		return d.Poll()
	}
}
