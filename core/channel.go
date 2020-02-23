package core

import (
	"gunplan.top/concurrentNet/buffer"
)

type Channel interface {
	Id() uint64
	Address() NetworkInet64
	Status() ConnectStatus
	AddTrigger(TimeTrigger)
	Type() ChannelType
	Parent() Channel
	buffer.Cached
}

type ParentChannel interface {
	Channel
	Listen(*NetworkInet64)
	Close()
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
	Write(buffer.ByteBuffer) error
	Read() (buffer.ByteBuffer, error)
	Close() error
	NetReset()
}

type channelImpl struct {
	id          uint64
	p           Channel
	address     NetworkInet64
	status      ConnectStatus
	channelType ChannelType
	Triggers    []TimeTrigger
	a           buffer.Allocator
	cc          ChannelCache
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

func (c *channelImpl) Parent() Channel {
	return c.p
}

func (c *channelImpl) Id() uint64 {
	return c.id
}

func (c *channelImpl) Release() {
	c.cc.release(c)
}

func (c *channelImpl) Reset() {
}

func (c *channelImpl) SetAlloc(a interface{}) {
	c.cc = a.(ChannelCache)
}
