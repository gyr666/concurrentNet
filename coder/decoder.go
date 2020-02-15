package coder

import "gunplan.top/concurrentNet/buffer"

type Decode func(buffer.ByteBuffer) interface{}
