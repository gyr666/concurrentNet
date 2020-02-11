package core

var Factory ChannelFactory = ChannelFactory{}

func NewConcurrentNet() Server{
	s := ServerImpl{}
	return s.Init()
}

type ChannelFactory struct {

}


func (c *ChannelFactory) NewChildChannelInstance() ChildChannel {
	return &childChannelImpl{}
}
func (c *ChannelFactory) NewParentChannelInstance() ParentChannel {
	return &parentChannelImpl{}
}
