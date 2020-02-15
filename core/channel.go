package core

import (
	"gunplan.top/concurrentNet/buffer"
)

type Channel interface {
	Address() NetworkInet64
	Status() ConnectStatus
	Write(buffer.ByteBuffer)
	Read() buffer.ByteBuffer
	Close() error
	Reset() error
	AddTrigger(TimeTrigger)
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
	readEvent(buffer.ByteBuffer)
	writeEvent() buffer.ByteBuffer
	closeEvent()
	exception(Throwable)
}
type channelImpl struct {
	id      uint64
	p       Channel
	address NetworkInet64
	status  ConnectStatus
	t       ChannelType
	Time    []TimeTrigger
	a		buffer.Allocator
	fd      int
}

func (c *channelImpl) Address() NetworkInet64 {
	return c.address
}

func (c *channelImpl) AddTrigger(t TimeTrigger) {
	c.Time = append(c.Time, t)
}

func (c *channelImpl) Status() ConnectStatus {
	return c.status
}
func (c *channelImpl) Write(buffer.ByteBuffer) {

}
func (c *channelImpl) Read() buffer.ByteBuffer {
	return nil
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



