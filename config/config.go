package config

import "gunplan.top/concurrentNet/buffer"

type ConfigDecoder func(byteBuffer buffer.ByteBuffer, config *Config) error
type Fetcher func() buffer.ByteBuffer

type GetConfig struct {
	c ConfigDecoder
	f Fetcher
}

func (g *GetConfig) Init(c ConfigDecoder, f Fetcher) {
	g.f = f
	g.c = c
}

func (g *GetConfig) Get() *Config {
	c := Config{}
	g.c(g.f(), &c)
	return &c
}
