package core

type Channel interface{
	Address() NetworkInet64
	Status() ConnectStatus
	Write(Data)
	Read() Data
	Close() error
	Reset() error
	Type()  ChannelType
	parent() Channel
}

type ChannelImpl struct{
	id	uint64
	p	Channel
	address	NetworkInet64
	status	ConnectStatus
	t	ChannelType
	fd	int
}

func (c *ChannelImpl) Address() NetworkInet64{
	return c.address
}

func (c *ChannelImpl) Status() ConnectStatus{
	return c.status
}
func (c *ChannelImpl) Write(Data){

}
func (c *ChannelImpl) Read() Data{
	return &dataImpl{}
}
func (c *ChannelImpl) Close() error{
	return nil
}
func (c *ChannelImpl) Reset() error{
	return nil
}
func (c *ChannelImpl) Type() ChannelType{
	return c.t
}
func (c *ChannelImpl) parent() Channel{
	return c.p
}

type ParentChannel struct{
	ChannelImpl
}

type ChildChannel struct{
	ChannelImpl
}
