package buffer

import (
	"errors"
	"gunplan.top/concurrentNet/util"
)

type OperatorMode int

const (
	INDEX_OUTOF_BOUND = "Index out of bound"

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
	AvailableReadSum() uint64
	Convert()
	FastMoveOut() []byte
	Mode() OperatorMode
	ShiftRN(n uint64) error
	ShiftWN(n uint64) error
}

type BaseByteBuffer struct {
	a       Allocator
	capital uint64
	WP      uint64
	RP      uint64
	RWType  bool
	s       []byte
}

func (b *BaseByteBuffer) Convert() {
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

func (b *BaseByteBuffer) ShiftRN(n uint64) error {
	if b.RP+n > b.capital {
		return errors.New(INDEX_OUTOF_BOUND)
	}
	b.RP += n
	return nil
}

func (b *BaseByteBuffer) ShiftWN(n uint64) error {
	if b.WP+n > b.capital {
		return errors.New(INDEX_OUTOF_BOUND)
	}
	b.WP += n
	return nil
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

func (b *BaseByteBuffer) AvailableReadSum() uint64 {
	return b.WP - b.RP - 1
}

func (b *BaseByteBuffer) FastMoveOut() []byte {
	return b.s
}
