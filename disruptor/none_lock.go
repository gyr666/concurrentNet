package disruptor

type Disruptor interface {
	Push(interface{})
	Poll() interface{}
}
