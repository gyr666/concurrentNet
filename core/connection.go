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

type Conn interface {
	buffer.Cached
	Event
	Id() uint64
	Address() NetworkInet64
	Status() ConnectStatus
	AddTrigger(TimeTrigger)
	Write(buffer.ByteBuffer) error
	Read() (buffer.ByteBuffer, error)
	Close() error
	NetReset()
}

func  NewConn() Conn {
	return &connImpl{}
}


type connImpl struct {
	id          uint64
	address     NetworkInet64
	status      ConnectStatus
	Triggers    []TimeTrigger
	a           buffer.Allocator
	cc          ConnCache
	fd          int


	l        Pipeline
	pool     threading.ThreadPool
	inCache  buffer.ByteBuffer
	outCache chan buffer.ByteBuffer
	call     CallBackEvent
	alloc    buffer.Allocator
}

func (c *connImpl) readEvent() {
	c.pool.Execwp(c.AsyncReadAndExecutePipeline, c.call, c.outCache)
	c.call.ReadEventAsyncExecuteComplete(c.id)
}

func (c *connImpl) AsyncReadAndExecutePipeline(i ...interface{}) {
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

func (c *connImpl) writeEvent() {
	c.pool.Execwp(c.AsyncWrite, c.call)
	c.call.WriteEventAsyncExecuteComplete(c.id)
}

func (c *connImpl) AsyncWrite(i ...interface{}) {
	e := i[0].(CallBackEvent)
	for len(c.outCache) != 0 {
		if err := c.Write(<-c.outCache); err != nil {
			e.PreWriteException(c.id, err)
			return
		}
	}
	e.WriteEventComplete(c.id)
}

func (c *connImpl) closeEvent() {
	err := c.Close()
	c.call.OperatorException(c.id, err)
}

func (c *connImpl) exception(t Throwable) {
	if !t.isUserDefine() {
		c.closeEvent()
	} else {
		c.l.doException(t, c)
	}
}

func (c *connImpl) Write(buffer.ByteBuffer) error {
	// block write is ok!
	//TODO

	//@gyr666 to implement
	return nil
}

func (c *connImpl) Read() (buffer.ByteBuffer, error) {
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

func (c *connImpl) Close() error {
	//TODO
	//@gyr666 to implement
	return nil
}

func (c *connImpl) NetReset() {
	//TODO
	//@gyr666 to implement
	c.call.PeerReset(c.id)
}

func (c *connImpl) Address() NetworkInet64 {
	return c.address
}

func (c *connImpl) AddTrigger(t TimeTrigger) {
	c.Triggers = append(c.Triggers, t)
}

func (c *connImpl) Status() ConnectStatus {
	return c.status
}

func (c *connImpl) Id() uint64 {
	return c.id
}

func (c *connImpl) Release() {
	c.cc.release(c)
}

func (c *connImpl) Reset() {
}

func (c *connImpl) SetAlloc(a interface{}) {
	c.cc = a.(ConnCache)
}