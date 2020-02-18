package core

import (
	"gunplan.top/concurrentNet/buffer"
)

type parentChannelImpl struct {
	channelImpl
}

func (p *parentChannelImpl) Listen(*NetworkInet64) {

}

func (c *parentChannelImpl) loop(byteBuffer buffer.ByteBuffer) {

}
