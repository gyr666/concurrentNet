package coder

import "gunplan.top/concurrentNet/buffer"

type Encode func(i interface{},b buffer.ByteBuffer) error
