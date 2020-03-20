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

type Channel interface {
	buffer.Cached
	Event
	Id() uint64
	Address() NetworkInet64
	Status() ChannelStatus
	AddTrigger(TimeTrigger)
	Write(buffer.ByteBuffer) error
	Read() (buffer.ByteBuffer, error)
	Close() error
	NetReset()
}

func NewChannel() Channel {
	return &channelImpl{}
}

type channelImpl struct {
	id       uint64
	address  NetworkInet64
	status   ChannelStatus
	Triggers []TimeTrigger
	a        buffer.Allocator
	cc       ChannelCache
	fd       int
	l        Pipeline
	pool     threading.ThreadPool
	inCache  buffer.ByteBuffer
	outCache chan buffer.ByteBuffer
	call     CallBackEvent
	alloc    buffer.Allocator
}

func (c *channelImpl) readEvent() {
	c.pool.Execwp(c.AsyncReadAndExecutePipeline, c.call, c.outCache)
	c.call.ReadEventAsyncExecuteComplete(c.id)
}

func (c *channelImpl) AsyncReadAndExecutePipeline(i ...interface{}) {
	e := i[0].(CallBackEvent)
	b, err := c.Read()
	if err != nil {
		e.PreReadException(c.id, err)
		return
	}
	e.ReadEventComplete(c.id)
	exp, err := c.l.doPipeline(b)
	if err != nil {
		e.OperatorException(c.id, err)
		return
	}
	(i[1]).(chan buffer.ByteBuffer) <- exp

	e.ReadPipelineComplete(c.id)
}

func (c *channelImpl) writeEvent() {
	c.pool.Execwp(c.AsyncWrite, c.call)
	c.call.WriteEventAsyncExecuteComplete(c.id)
}

func (c *channelImpl) AsyncWrite(i ...interface{}) {
	e := i[0].(CallBackEvent)
	for len(c.outCache) != 0 {
		if err := c.Write(<-c.outCache); err != nil {
			e.PreWriteException(c.id, err)
			return
		}
	}
	e.WriteEventComplete(c.id)
}

func (c *channelImpl) closeEvent() {
	err := c.Close()
	c.call.OperatorException(c.id, err)
}

func (c *channelImpl) exception(t Throwable) {
	if !t.isUserDefine() {
		c.closeEvent()
	} else {
		c.l.doException(t, c)
	}
}

func (c *channelImpl) Write(buffer.ByteBuffer) error {
	// block write is ok!
	//TODO

	//@gyr666 to implement
	return nil
}

func (c *channelImpl) Read() (buffer.ByteBuffer, error) {
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
		bytes := (*c.inCache.FastMoveOut())[head : head+sz]

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
		bytes = *buf.FastMoveOut()
		n, err = unix.Read(c.fd, bytes[sz:])
		if n == 0 || err != nil {
			if err == unix.EAGAIN {
				return c.inCache, nil
			}
			return c.alloc.Alloc(0), errChannelClose
		}
		//there is never index out of bound
		_ = buf.ShiftWN(uint64(n))
		c.inCache = buf
		initSz <<= 1
	}
}

func (c *channelImpl) Close() error {
	//TODO
	//@gyr666 to implement
	return nil
}

func (c *channelImpl) NetReset() {
	//TODO
	//@gyr666 to implement
	c.call.PeerReset(c.id)
}

func (c *channelImpl) Address() NetworkInet64 {
	return c.address
}

func (c *channelImpl) AddTrigger(t TimeTrigger) {
	c.Triggers = append(c.Triggers, t)
}

func (c *channelImpl) Status() ChannelStatus {
	return c.status
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
