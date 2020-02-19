package core

var Factory = ChannelFactory{}

func NewConcurrentNet() Server {
	s := ServerImpl{}
	return s.Init()
}

type ChannelFactory struct {
	cache ChannelCache
}

func (c *ChannelFactory) NewChildChannelInstance() ChildChannel {
	if c.cache != nil {
		return c.cache.Acquire(Child).(ChildChannel)
	}
	return &childChannelImpl{}
}

func (c *ChannelFactory) NewParentChannelInstance() ParentChannel {
	return c.cache.Acquire(Child).(ParentChannel)
}
