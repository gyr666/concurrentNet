package buffer

type Develop interface {
	Init(uint64) error
	Destroy() error
}

type ByteBufferDevelop interface {
	Init(uint64, Allocator)
	Destroy() error
}
type IOer interface {
	Writer
	Reader
}

type Writer interface {
	Write([]byte) error
}
type Reader interface {
	Read(uint64) ([]byte, error)
}
