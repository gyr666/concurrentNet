package config

import (
	"fmt"
	"io/ioutil"

	"gunplan.top/concurrentNet/buffer"
)

type GetFromDefaultStrategy struct {
}

func (g *GetFromDefaultStrategy) Fill(c *Config) error {
	return nil
}

type GetConfigFromFile struct {
	filePath string
}

func (g *GetConfigFromFile) FetchToMemory() buffer.ByteBuffer {
	f, err := ioutil.ReadFile(g.filePath)
	bf := buffer.NewSandBufferAllocator().Alloc(uint64(len(f)))
	if err != nil {
		fmt.Println("read fail", err)
	}
	bf.Write(f)
	return bf
}

func (g *GetConfigFromFile) Decode() ConfigDecoder {
	return nil
}
