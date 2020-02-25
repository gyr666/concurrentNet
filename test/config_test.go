package test

import (
	"io/ioutil"
	"testing"

	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/config"
)

func TestGetConfig(t *testing.T) {
	c := config.GetConfig{}
	c.Init(config.LineDecoder, func() buffer.ByteBuffer {
		data, _ := ioutil.ReadFile("../params/boot.conf")
		b := buffer.NewLikedBufferAllocator().Alloc(uint64(len(data)))
		b.Write(data)
		return b
	})
	c.Get()

}
