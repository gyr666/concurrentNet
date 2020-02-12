package buffer

import (
	"gunplan.top/concurrentNet/util"
)

type OperatorMode int

const (
	READ  OperatorMode = 1
	WRITE OperatorMode = 1 << 1
)

func getMode(b bool) OperatorMode {
	if b {
		return READ
	}
	return WRITE
}

type Allocator interface {
	Develop
	Alloc(uint64) ByteBuffer
	release(ByteBuffer)
	PoolSize() uint64
}

type ByteBuffer interface {
	Develop
	IOer
	Release()
	Size() uint64
	convert()
	Mode() OperatorMode
}

type BaseByteBuffer struct {
	a       Allocator
	capital uint64
	WP      uint64
	RP      uint64
	RWType  bool
	s 		[]byte
}

func (b *BaseByteBuffer) convert() {
	b.RWType = !b.RWType
}

func (b *BaseByteBuffer) Mode() OperatorMode {
	return getMode(b.RWType)
}

func (b *BaseByteBuffer) Size() uint64 {
	return b.capital
}

func (b *BaseByteBuffer) Read(i int) ([]byte, error) {
	return util.StandRead(i,b.s,b.capital,&b.RP)
}

func (b *BaseByteBuffer) Write(_b []byte) error {
	return util.StandWrite(b.s,b.capital,&b.WP,_b)
}

func (a *BaseByteBuffer) Init(s uint64) error {
	a.capital = s
	a.s = make([]byte, a.capital)
	return nil
}

func (a *BaseByteBuffer) Destroy() error {
	a.reset()
	a.s = nil
	//help gc
	return nil
}

func (b *BaseByteBuffer) Release() {
	b.a.release(b)
}

func (b *BaseByteBuffer) reset(){
	b.WP = 0
	b.RP = 0
}