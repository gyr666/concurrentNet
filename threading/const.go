package threading

type PoolState uint8
type ControlType uint8

const (
	STOPPED  PoolState = 0
	RUNNING  PoolState = 1 << 1
	WAITING  PoolState = 1 << 2
	STOPPING PoolState = 1 << 3
	BOOTING  PoolState = 1 << 4
)

const (
	ShutdownNow ControlType = 1 << 1
	SHUTDOWN    ControlType = 1 << 2
)
