package threading

type PoolState uint8
type ControlType uint8

const (
	STOP        PoolState = 1
	RUNNING     PoolState = 1 << 1
	WAITING     PoolState = 1 << 2
	SHUTDOWNING PoolState = 1 << 3
)

const (
	STOPANY  ControlType = 1
	STOPALL  ControlType = 1 << 1
	SHUTDOWN ControlType = 1 << 2
)
