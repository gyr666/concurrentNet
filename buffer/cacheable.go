package buffer

type Cached interface {
	Release()
	Reset()
	SetAlloc(interface{})
}
