package buffer

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
	WP      uint
	RP      uint
	RWType  bool
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
