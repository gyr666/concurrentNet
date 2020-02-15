package core

import (
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/threading"
)

type childChannelImpl struct {
	channelImpl
	l Pipeline
	pool threading.ThreadPool
	cache chan buffer.ByteBuffer
}

func (c *childChannelImpl) readEvent(){
	c.pool.Execwp(func(i ...interface{}) {
		(i[0]).(chan buffer.ByteBuffer) <- c.l.doPipeline((c.Read()))
	},c.cache)
}

func (c *childChannelImpl) writeEvent() {
	for ;len(c.cache)!=0; {
		c.Write(<-c.cache)
	}
}

func (c *childChannelImpl) closeEvent(){
	_ = c.Close()
}

func (c *childChannelImpl) exception(t Throwable){
	if !t.isUserDefine() {
		c.closeEvent()
	} else {
		c.l.doException(t,c)
	}
}

func (c *channelImpl) Write(buffer.ByteBuffer) {
	// block write is ok!
	//TODO
	//@gyr666 to implement
}

func (c *channelImpl) Read() buffer.ByteBuffer {
	// block read is ok!
	//TODO
	//@gyr666 to implement
	return nil
}

func (c *channelImpl) Close() error {
	//TODO
	//@gyr666 to implement
	return nil
}

func (c *channelImpl) Reset() error {
	//TODO
	//@gyr666 to implement
	return nil
}