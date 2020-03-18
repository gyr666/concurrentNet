package config

import "gunplan.top/concurrentNet/buffer"

type GetFromFileStrategy struct {
}

func DefaultFetcher() buffer.ByteBuffer {
	return nil
}
