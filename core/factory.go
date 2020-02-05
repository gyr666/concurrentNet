package core

var Factory ChannelFactory = ChannelFactory{}

func NewConcurrentNet() Server{
	s := ServerImpl{}
	return s.Init()
}

type ChannelFactory struct {

}


func (c *ChannelFactory) NewChildChannelInstance() ChildChannel {
	return ChildChannel{}
}
func (c *ChannelFactory) NewParentChannelInstance() ParentChannel {
	return ParentChannel{}
}
