package event_loop

type EventType uint
type ExecCode uint

const (
	SYNC   ExecCode = 1
	ASYNC  ExecCode = 1 << 1
	ERROR  ExecCode = 1 << 2
	HALT   ExecCode = 1 << 4
	NORMAL ExecCode = 1 << 3
)

const (
	CONNECTION EventType = 1
	READ       EventType = 1 << 1
	WRITE      EventType = 1 << 2
	CLOSE      EventType = 1 << 3
	RESET      EventType = 1 << 4
)

type BaseEventLoop interface {
	StartLoop() ExecCode
	haltLoop() ExecCode
	resumeLoop() ExecCode
	stopLoop() ExecCode
}

type EventLoop interface {
	BaseEventLoop
	processEvent(Event) error
}

type Event interface {
	EType() EventType
	Attach() interface{}
}
