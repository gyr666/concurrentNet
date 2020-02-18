package core

import (
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/threading"
)

type childChannelImpl struct {
	channelImpl
	l     Pipeline
	pool  threading.ThreadPool
	cache chan buffer.ByteBuffer
	call  CallBackEvent
}

func (c *childChannelImpl) readEvent() {
	c.pool.Execwp(c.AsyncReadAndExecutePipeline, c.call, c.cache)
	c.call.ChannelReadEventAsyncExecuteComplete(c.id)
}

func (c *childChannelImpl) AsyncReadAndExecutePipeline(i ...interface{}) {
	e := i[0].(CallBackEvent)
	b, err := c.ReadAll()
	if err != nil {
		e.ChannelPreReadException(c.id, err)
		return
	}
	e.ChannelReadEventComplete(c.id)
	for index := range b {
		exp, err := c.l.doPipeline(b[index])
		if err != nil {
			e.ChannelOperatorException(c.id, err)
			continue
		}
		(i[1]).(chan buffer.ByteBuffer) <- exp
	}
	e.ChannelReadPipelineComplete(c.id)
}

func (c *childChannelImpl) writeEvent() {
	c.pool.Execwp(c.AsyncWrite, c.call)
	c.call.ChannelWriteEventAsyncExecuteComplete(c.id)
}

func (c *childChannelImpl) AsyncWrite(i ...interface{}) {
	e := i[0].(CallBackEvent)
	for len(c.cache) != 0 {
		if err := c.Write(<-c.cache); err != nil {
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
	//TODO
	//@gyr666 to implement
	return nil, nil
}

func (c *childChannelImpl) Close() error {
	//TODO
	//@gyr666 to implement
	return nil
}

func (c *childChannelImpl) Reset() {
	//TODO
	//@gyr666 to implement
	c.call.ChannelPeerReset(c.id)
}

func (c *childChannelImpl) ReadAll() ([]buffer.ByteBuffer, error) {
	return nil, nil
}
