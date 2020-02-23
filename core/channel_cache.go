package core

import (
	"sync"
)

type ChannelCache interface {
	Acquire(ChannelType) Channel
	release(c Channel)
}

func NewChannelCache() ChannelCache {
	c := channelCacheImpl{childLoad: 100}
	return c.Init()
}

type LinkedChannel struct {
	c    Channel
	next *LinkedChannel
}

type channelCacheMeta struct {
	first *LinkedChannel
	l     sync.Mutex
	t     ChannelType
	index uint8
	a     ChannelCache
	load  uint8
}

func (c *channelCacheMeta) init(channelType ChannelType, load uint8, a ChannelCache, index uint8) {
	c.t = channelType
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
		if c.t == Child {
			lc[i].c = &childChannelImpl{}
			lc[i].c.(*childChannelImpl).id = uint64(i)
		} else {
			lc[i].c = &parentChannelImpl{}
		}
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
	childMeta  []channelCacheMeta
	parentMeta channelCacheMeta
	childLoad  uint8
	i          uint8
}

func (c *channelCacheImpl) Init() ChannelCache {
	c.parentMeta = channelCacheMeta{t: Parent, load: 200, a: c, index: 0}
	c.childMeta = make([]channelCacheMeta, c.childLoad)
	for i := range c.childMeta {
		c.childMeta[i].init(Child, 200, c, uint8(i))
	}
	return c
}

func (c *channelCacheImpl) Acquire(ct ChannelType) Channel {
	if ct == Parent {
		return c.parentMeta.Alloc()
	}
	for {
		in := c.next()
		if c.childMeta[in].checkSelfAddUpdate() {
			return c.childMeta[in].Alloc()
		}
	}
}

func (c *channelCacheImpl) release(ch Channel) {
	if ch.Type() == Parent {
		c.parentMeta.release(ch)
		return
	}
	c.childMeta[c.next()].release(ch)
}

func (c *channelCacheImpl) next() uint8 {
	c.i++
	if c.i == c.childLoad {
		c.i = 0
	}
	return c.i
}
