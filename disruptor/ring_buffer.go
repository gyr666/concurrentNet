package disruptor

import (
	"sync/atomic"
)

type disruptorImpl struct {
	buffer []interface{}
	p      uint64
	c      uint64
}

func (d *disruptorImpl) Init(n int) {
	d.buffer = make([]interface{}, n)
	d.p = 0
	d.c = 0
}

func (d *disruptorImpl) Push(i interface{}) {
	d.buffer[d.p] = i
	d.p++
}

func (d *disruptorImpl) Poll() interface{} {
	v := d.buffer[d.c]
	if atomic.CompareAndSwapUint64(&d.p, d.p, d.p-2) {
		return v
	} else {
		return d.Poll()
	}
}
