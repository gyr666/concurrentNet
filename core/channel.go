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

}
func (c *ChannelImpl) Address() NetworkInet64{
	return NetworkInet64{Address:"[::]:0"}
}
func (c *ChannelImpl) Status() ConnectStatus{
	return NORMAL
}
func (c *ChannelImpl) Write(Data){

}
func (c *ChannelImpl) Read() Data{
	return &dataImpl{}
}
func (c *ChannelImpl) Close() error{
	return nil
}
func (c *ChannelImpl)Reset() error{
	return nil
}
func (c *ChannelImpl)Type() ChannelType{
	return ChildChannel_
}
func (c *ChannelImpl)parent() Channel{
	return nil
}

type ParentChannel struct{
	ChannelImpl
}

type ChildChannel struct{
	ChannelImpl
}
