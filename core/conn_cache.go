package core

import (
	"sync"
)

type ConnCache interface {
	Acquire() Conn
	release(c Conn)
}

func NewChannelCache() ConnCache {
	c := connCacheImpl{load: 100}
	return c.Init()
}

type LinkedChannel struct {
	c    Conn
	next *LinkedChannel
}

type channelCacheMeta struct {
	first *LinkedChannel
	l     sync.Mutex
	index uint8
	a     ConnCache
	load  uint8
}

func (c *channelCacheMeta) init(load uint8, a ConnCache, index uint8) {
	c.load = load
	c.a = a
	c.index = index
	c.createList()
}

func (c *channelCacheMeta) createList() {
	c.l.Lock()
	defer c.l.Unlock()
	lc := make([]LinkedChannel, c.load)
	for i := 0; i < len(lc)-1; i++ {
		lc[i].c = &connImpl{}
		lc[i].c.(*connImpl).id = uint64(i)
		lc[i].c.SetAlloc(c.a)
		lc[i].next = &lc[i+1]
	}
	c.first = &lc[0]
}

func (c *channelCacheMeta) Alloc() Conn {
	c.l.Lock()
	defer c.l.Unlock()
	v := c.first
	c.first = v.next
	v.next = nil
	return v.c
}

func (c *channelCacheMeta) release(ch Conn) {
	c.l.Lock()
	defer c.l.Unlock()
	ch.Reset()
	w := &LinkedChannel{c: ch}
	v := c.first
	c.first = w
	w.next = v
}

func (c *channelCacheMeta) checkSelfAddUpdate() bool {
	if c.first == nil {
		go c.createList()
		return false
	}
	return true
}

type connCacheImpl struct {
	meta []channelCacheMeta
	load uint8
	i    uint8
}

func (c *connCacheImpl) Init() ConnCache {
	c.meta = make([]channelCacheMeta, c.load)
	for i := range c.meta {
		c.meta[i].init( 200, c, uint8(i))
	}
	return c
}

func (c *connCacheImpl) Acquire() Conn {
	for {
		in := c.next()
		if c.meta[in].checkSelfAddUpdate() {
			return c.meta[in].Alloc()
		}
	}
}

func (c *connCacheImpl) release(ch Conn) {
	c.meta[c.next()].release(ch)
}

func (c *connCacheImpl) next() uint8 {
	c.i++
	if c.i == c.load {
		c.i = 0
	}
	return c.i
}
