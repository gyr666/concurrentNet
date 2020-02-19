package core

import (
	"sync"
)

type ChannelCache interface {
	Acquire(ChannelType) Channel
	release(c Channel)
}

func NewChannelCache() ChannelCache {
	c := channelCacheImpl{load: 20}
	return c.Init()
}

type LinkedChannel struct {
	c    Channel
	next *LinkedChannel
}

type channelCacheImpl struct {
	p    *LinkedChannel
	c    *LinkedChannel
	cl   sync.Mutex
	pl   sync.Mutex
	load uint8
}

func (c *channelCacheImpl) Init() ChannelCache {
	c.createList(Child)
	c.createList(Parent)
	return c
}

func (c *channelCacheImpl) createList(channelType ChannelType) *LinkedChannel {
	lc := make([]LinkedChannel, c.load)
	for i := 0; i < len(lc)-1; i++ {
		if channelType == Child {
			lc[i].c = &childChannelImpl{}
		} else {
			lc[i].c = &parentChannelImpl{}
		}
		lc[i].c.SetAlloc(c)
		lc[i].next = &lc[i+1]
	}
	ls, _ := c.findList(channelType)
	*ls = &lc[0]
	return &lc[0]
}

func (c *channelCacheImpl) head(ct ChannelType) Channel {
	lc, l := c.findList(ct)
	l.Lock()
	defer l.Unlock()
	v := *lc
	*lc = v.next
	return v.c
}

func (c *channelCacheImpl) Acquire(ct ChannelType) Channel {
	return c.head(ct)
}

func (c *channelCacheImpl) release(ch Channel) {
	ch.Reset()
	lc, l := c.findList(ch.Type())
	l.Lock()
	l.Unlock()
	if *lc == nil {
		c.createList(ch.Type())
	}
	v := *lc
	*lc = &LinkedChannel{c: ch}
	(*lc).next = v
}

func (c *channelCacheImpl) findList(t ChannelType) (**LinkedChannel, *sync.Mutex) {
	if t == Child {
		return &c.c, &c.cl
	}
	return &c.p, &c.pl
}
