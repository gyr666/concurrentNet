package buffer

type Develop interface {
	Init() error
	Destroy() error
}

type IOer interface {
	Writer
	Reader
}

type Writer interface {
	write([]byte) error
}
type Reader interface {
	readInt() (int, error)
	readNByte(int) ([]byte, error)
}
