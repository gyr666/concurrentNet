package coder

type Encode func(interface{}) []byte
type Decode func([]byte) interface{}

