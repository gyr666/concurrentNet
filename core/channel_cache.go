package core

import (
	"sync"
)

type ChannelCache interface {
	Acquire() Channel
	release(c Channel)
}

func NewChannelCache() ChannelCache {
	c := channelCacheImpl{load: 100}
	return c.Init()
}

type LinkedChannel struct {
	c    Channel
	next *LinkedChannel
}

type channelCacheMeta struct {
	first *LinkedChannel
	l     sync.Mutex
	index uint8
	a     ChannelCache
	load  uint8
}

func (c *channelCacheMeta) init(load uint8, a ChannelCache, index uint8) {
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
		lc[i].c = &channelImpl{}
		lc[i].c.(*channelImpl).id = uint64(i)
		lc[i].c.SetAlloc(c.a)
		lc[i].next = &lc[i+1]
	}
	c.first = &lc[0]
}

func (c *channelCacheMeta) Alloc() Channel {
	c.l.Lock()
	defer c.l.Unlock()
	v := c.first
	c.first = v.next
	v.next = nil
	return v.c
}

func (c *channelCacheMeta) release(ch Channel) {
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

type channelCacheImpl struct {
	meta []channelCacheMeta
	load uint8
	i    uint8
}

func (c *channelCacheImpl) Init() ChannelCache {
	c.meta = make([]channelCacheMeta, c.load)
	for i := range c.meta {
		c.meta[i].init(200, c, uint8(i))
	}
	return c
}

func (c *channelCacheImpl) Acquire() Channel {
	for {
		in := c.next()
		if c.meta[in].checkSelfAddUpdate() {
			return c.meta[in].Alloc()
		}
	}
}

func (c *channelCacheImpl) release(ch Channel) {
	c.meta[c.next()].release(ch)
}

func (c *channelCacheImpl) next() uint8 {
	c.i++
	if c.i == c.load {
		c.i = 0
	}
	return c.i
}
