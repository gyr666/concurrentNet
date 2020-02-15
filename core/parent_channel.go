package core

import (
	"gunplan.top/concurrentNet/buffer"
)

type parentChannelImpl struct {
	channelImpl
}


func (p *parentChannelImpl) Listen(*NetworkInet64) {

}
func (c *childChannelImpl) loop(byteBuffer buffer.ByteBuffer) {

}
