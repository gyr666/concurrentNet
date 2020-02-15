package core

import (
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/threading"
)

type parentChannelImpl struct {
	channelImpl
}


func (p *parentChannelImpl) Listen(*NetworkInet64) {
}
func (c *childChannelImpl) loop(pool threading.ThreadPool,byteBuffer buffer.ByteBuffer) {

}
