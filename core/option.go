package core

type OptionType interface {
	doSet(interface{}) error
}

type BackLog struct {
}

type BufferLength struct {
}

type NetWorkType struct {
}
type StreamType struct {
}

func (b *BackLog) doSet(interface{}) error      { return nil }
func (b *BufferLength) doSet(interface{}) error { return nil }
func (b *NetWorkType) doSet(interface{}) error  { return nil }
func (b *StreamType) doSet(interface{}) error   { return nil }
