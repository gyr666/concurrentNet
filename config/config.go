package config

import "gunplan.top/concurrentNet/buffer"

type ConfigDecoder func(byteBuffer buffer.ByteBuffer, config *Config) error
type Fetcher func() buffer.ByteBuffer

type GetConfig struct {
	C ConfigDecoder
	F Fetcher
}

func (g *GetConfig) Get() *Config {
	c := Config{}
	g.C(g.F(), &c)
	return &c
}
