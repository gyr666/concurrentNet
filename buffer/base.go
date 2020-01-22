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
	Write([]byte) error
}
type Reader interface {
	Read(int) ([]byte, error)
}
