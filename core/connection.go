package core

import (
	"golang.org/x/sys/unix"
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/threading"
)



type Event interface {
	readEvent()
	writeEvent()
	closeEvent()
	exception(Throwable)
}

type ChildChannel interface {
	Event
	Id() uint64
	Address() NetworkInet64
	Status() ConnectStatus
	AddTrigger(TimeTrigger)
	Type() ChannelType
	Write(buffer.ByteBuffer) error
	Read() (buffer.ByteBuffer, error)
	Close() error
	NetReset()
}

type channelImpl struct {
	id          uint64
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
type childChannelImpl struct {
	channelImpl
	l        Pipeline
	pool     threading.ThreadPool
	inCache  buffer.ByteBuffer
	outCache chan buffer.ByteBuffer
	call     CallBackEvent
	alloc    buffer.Allocator
}

func (c *childChannelImpl) readEvent() {
	c.pool.Execwp(c.AsyncReadAndExecutePipeline, c.call, c.outCache)
	c.call.ChannelReadEventAsyncExecuteComplete(c.id)
}

func (c *childChannelImpl) AsyncReadAndExecutePipeline(i ...interface{}) {
	e := i[0].(CallBackEvent)
	b, err := c.Read()
	if err != nil {
		e.ChannelPreReadException(c.id, err)
		return
	}
	e.ChannelReadEventComplete(c.id)
	exp, err := c.l.doPipeline(b)
	if err != nil {
		e.ChannelOperatorException(c.id, err)
		return
	}
	(i[1]).(chan buffer.ByteBuffer) <- exp

	e.ChannelReadPipelineComplete(c.id)
}

func (c *childChannelImpl) writeEvent() {
	c.pool.Execwp(c.AsyncWrite, c.call)
	c.call.ChannelWriteEventAsyncExecuteComplete(c.id)
}

func (c *childChannelImpl) AsyncWrite(i ...interface{}) {
	e := i[0].(CallBackEvent)
	for len(c.outCache) != 0 {
		if err := c.Write(<-c.outCache); err != nil {
			e.ChannelPreWriteException(c.id, err)
			return
		}
	}
	e.ChannelWriteEventComplete(c.id)
}

func (c *childChannelImpl) closeEvent() {
	err := c.Close()
	c.call.ChannelOperatorException(c.id, err)
}

func (c *childChannelImpl) exception(t Throwable) {
	if !t.isUserDefine() {
		c.closeEvent()
	} else {
		c.l.doException(t, c)
	}
}

func (c *childChannelImpl) Write(buffer.ByteBuffer) error {
	// block write is ok!
	//TODO

	//@gyr666 to implement
	return nil
}

func (c *childChannelImpl) Read() (buffer.ByteBuffer, error) {
	// block read is ok!
	var (
		err error
		n   int
	)
	defer func() {
		if n == 0 || err != nil {
			if err != unix.EAGAIN {
				c.closeEvent()
			}
		}
	}()
	initSz := uint64(1024)

	for {
		//之前未读完数据的切片
		head := c.inCache.GetRP()
		sz := c.inCache.AvailableReadSum()
		bytes := c.inCache.FastMoveOut()[head : head+sz]

		//申请新的 buffer
		buf := c.alloc.Alloc(initSz)
		//把之前数据的写入新的 buffer
		if sz != 0 {
			err = buf.Write(bytes)
			if err != nil {
				return c.alloc.Alloc(0), nil
			}
		}
		//再将连接读到的数据写入新的 buffer
		bytes = buf.FastMoveOut()
		n, err = unix.Read(c.fd, bytes[sz:])
		if n == 0 || err != nil {
			if err == unix.EAGAIN {
				return c.inCache, nil
			}
			return c.alloc.Alloc(0), errConnClose
		}
		//there is never index out of bound
		_ = buf.ShiftWN(uint64(n))
		c.inCache = buf
		initSz <<= 1
	}
}

func (c *childChannelImpl) Close() error {
	//TODO
	//@gyr666 to implement
	return nil
}

func (c *childChannelImpl) NetReset() {
	//TODO
	//@gyr666 to implement
	c.call.ChannelPeerReset(c.id)
}
