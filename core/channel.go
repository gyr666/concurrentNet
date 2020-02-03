package core

type Channel interface {
	Address() NetworkInet64
	Status() ConnectStatus
	Write(Data)
	Read() Data
	Close() error
	Reset() error
	AddTigger(TimeTigger)
	Type() ChannelType
	parent() Channel
}

type ParentChannel interface {
	Channel
	Listen(*NetworkInet64)
	loop()
}

type ChildChannel interface {
	Channel
}
type channelImpl struct {
	id      uint64
	p       Channel
	address NetworkInet64
	status  ConnectStatus
	t       ChannelType
	fd      int
}

func (c *channelImpl) Address() NetworkInet64 {
	return c.address
}

func (c *channelImpl) AddTigger(t TimeTigger) {
}

func (c *channelImpl) Status() ConnectStatus {
	return c.status
}
func (c *channelImpl) Write(Data) {

}
func (c *channelImpl) Read() Data {
	return &dataImpl{}
}
func (c *channelImpl) Close() error {
	return nil
}
func (c *channelImpl) Reset() error {
	return nil
}
func (c *channelImpl) Type() ChannelType {
	return c.t
}
func (c *channelImpl) parent() Channel {
	return c.p
}

type parentChannelImpl struct {
	channelImpl
}

type childChannelImpl struct {
	channelImpl
}

func (p *parentChannelImpl) Listen(*NetworkInet64) {
}
func (p *parentChannelImpl) loop() {
}
