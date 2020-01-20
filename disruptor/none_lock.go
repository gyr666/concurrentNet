package disruptor

type Disruptor interface {
	Push(interface{})
	Poll() interface{}
}

type Sequence interface {
	Next() (uint64, bool)
}
