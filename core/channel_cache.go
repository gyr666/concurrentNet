package core

type ChannelCache interface {
	Acquire(ChannelType) Channel
	release(c Channel)
}

type LinkedChannel struct {
	c    Channel
	next *LinkedChannel
}

type channelCacheImpl struct {
	p *LinkedChannel
	c *LinkedChannel
}

func (c *channelCacheImpl) Init() {
	c.c = c.createList(Child)
	c.p = c.createList(Parent)
}

func (c *channelCacheImpl) createList(channelType ChannelType) *LinkedChannel {
	lc := make([]LinkedChannel, 20)
	for i := 0; i < len(lc)-1; i++ {
		if channelType == Child {
			lc[i].c = &childChannelImpl{}
		} else {
			lc[i].c = &parentChannelImpl{}
		}
		lc[i].next = &lc[i+1]
	}
	return &lc[0]
}

func (c *channelCacheImpl) head(ct ChannelType) Channel {
	var lc = c.findList(ct)
	v := *lc
	(*lc) = v.next
	return v.c
}

func (c *channelCacheImpl) Acquire(ct ChannelType) Channel {
	return c.head(ct)
}

func (c *channelCacheImpl) release(ch Channel) {
	ch.Reset()
	lc := c.findList(ch.Type())
	v := *lc
	*lc = &LinkedChannel{c: ch}
	(*lc).next = v
}

func (c *channelCacheImpl) findList(t ChannelType) **LinkedChannel {
	if t == Child {
		return &c.c
	}
	return &c.p
}
