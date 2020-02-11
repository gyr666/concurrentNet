package buffer

type Develop interface {
	Init(s uint64) error
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
	Read(int) ([]byte, error)
}
