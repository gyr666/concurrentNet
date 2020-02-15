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

func (c *childChannelImpl) readEvent(b buffer.ByteBuffer){
	c.pool.Execwp(func(i ...interface{}) {
		(i[1]).(chan buffer.ByteBuffer) <- c.l.doPipeline(i[0].(buffer.ByteBuffer))
	},b,c.cache)
}

func (c *childChannelImpl) writeEvent() buffer.ByteBuffer{
	return <- c.cache
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
