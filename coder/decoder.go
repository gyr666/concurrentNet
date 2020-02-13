package coder

import "gunplan.top/concurrentNet/buffer"

type decode func(buffer.ByteBuffer) interface{}
