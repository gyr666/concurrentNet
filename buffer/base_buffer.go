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
	OperatorTimes() uint64
	release(ByteBuffer)
	PoolSize() uint64
	AllocSize() uint64
}

type ByteBuffer interface {
	ByteBufferDevelop
	IOer
	Cached
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
	s       []byte
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

func (b *BaseByteBuffer) Read(i uint64) ([]byte, error) {
	return util.StandRead(i, b.s, b.capital, &b.RP)
}

func (b *BaseByteBuffer) ReadAll() ([]byte, error) {
	return util.StandRead((b.RP), b.s, b.capital, &b.RP)
}

func (b *BaseByteBuffer) Write(_b []byte) error {
	return util.StandWrite(b.s, b.capital, &b.WP, _b)
}

func (b *BaseByteBuffer) Init(s uint64, all Allocator) {
	b.capital = s
	b.a = all
	b.s = make([]byte, b.capital)
}

func (b *BaseByteBuffer) Destroy() error {
	b.Reset()
	b.s = nil
	//help gc
	return nil
}

func (b *BaseByteBuffer) Release() {
	b.a.release(b)
}

func (b *BaseByteBuffer) Reset() {
	b.WP = 0
	b.RP = 0
}

func (b *BaseByteBuffer) SetAlloc(i interface{}) {
	b.a = i.(Allocator)
}
