package config

import (
	"gunplan.top/concurrentNet/buffer"
	"io/ioutil"
)

type GetFromFileStrategy struct {
}

func DefaultFetcher() buffer.ByteBuffer {
	a := buffer.NewLikedBufferAllocator()
	data, _ := ioutil.ReadFile("params/boot.conf")
	b := a.Alloc(1)
	b.FastMoveIn(&data)
	return b
}
