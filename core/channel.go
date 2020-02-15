package core

import (
	"gunplan.top/concurrentNet/buffer"
)

type Channel interface {
	Address() NetworkInet64
	Status() ConnectStatus
	AddTrigger(TimeTrigger)
	Type() ChannelType
	parent() Channel
}

type ParentChannel interface {
	Channel
	Listen(*NetworkInet64)
	loop()
}

type Event interface {
	readEvent()
	writeEvent()
	closeEvent()
	exception(Throwable)
}

type ChildChannel interface {
	Channel
	Event
	Write(buffer.ByteBuffer)
	Read() buffer.ByteBuffer
	Close() error
	Reset() error
}

type channelImpl struct {
	id          uint64
	p           Channel
	address     NetworkInet64
	status      ConnectStatus
	channelType ChannelType
	Triggers    []TimeTrigger
	a           buffer.Allocator
	fd          int
}

func (c *channelImpl) Address() NetworkInet64 {
	return c.address
}

func (c *channelImpl) AddTrigger(t TimeTrigger) {
	c.Triggers = append(c.Triggers, t)
}

func (c *channelImpl) Status() ConnectStatus {
	return c.status
}

func (c *channelImpl) Type() ChannelType {
	return c.channelType
}

func (c *channelImpl) parent() Channel {
	return c.p
}



