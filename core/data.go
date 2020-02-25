package core

type TransferEvent interface {
}

type Data interface {
	ConsumeEvent() TransferEvent
	ProduceEvent(TransferEvent)
	GetThrow() *Throwable
	SetThrow(*Throwable)
	Channel() Channel
	GetData() interface{}
	SetData(interface{})
}

type dataImpl struct {
	data  interface{}
	event chan TransferEvent
	throw *Throwable
	c     Channel
}

func (d *dataImpl) Channel() Channel {
	return d.c
}

func (d *dataImpl) ConsumeEvent() TransferEvent {
	return d.event
}

func (d *dataImpl) ProduceEvent(t TransferEvent) {
	d.event <- t
}

func (d *dataImpl) GetThrow() *Throwable {
	return d.throw
}

func (d *dataImpl) SetThrow(t *Throwable) {
	d.throw = t
}

func (d *dataImpl) GetData() interface{} {
	return d.data
}

func (d *dataImpl) SetData(data interface{}) {
	d.data = data
}
