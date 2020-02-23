package config

import "gunplan.top/concurrentNet/buffer"

type ConfigStrategy interface {
	Fill(*Config) error
}

type ConfigDecoder interface {
	Decode(byteBuffer buffer.ByteBuffer, config *Config)
}

type ServiceGetter interface {
	FetchToMemory() buffer.ByteBuffer
	Decode() ConfigDecoder
}

type GetConfig struct {
	sg ServiceGetter
}

func (g *GetConfig) Get() *Config {
	b := g.sg.FetchToMemory()
	c := Config{}
	g.sg.Decode().Decode(b, &c)
	return &c
}
