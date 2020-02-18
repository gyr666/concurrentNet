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
	c.pool.Execwp(func(i ...interface{}) {
		b, err := c.Read()
		if err != nil {
			c.call.ChannelPreReadException(c.id, err)
		}
		(i[0]).(chan buffer.ByteBuffer) <- c.l.doPipeline(b)
		c.call.ChannelReadEventComplete(c.id)
	}, c.cache)
	c.call.ChannelReadEventAsyncExecuteComplete(c.id)
}

func (c *childChannelImpl) writeEvent() {
	c.pool.Execwp(func(i ...interface{}) {
		for len(c.cache) != 0 {
			if err := c.Write(<-c.cache); err != nil {
				i[0].(CallBackEvent).ChannelPreWriteException(c.id, err)
				return
			}
		}
		i[0].(CallBackEvent).ChannelWriteEventComplete(c.id)
	}, c.call)
	c.call.ChannelWriteEventAsyncExecuteComplete(c.id)
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
