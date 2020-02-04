package core

type WaitType uint8

const(
	SYNC  WaitType = 1
	ASYNC WaitType = 1<<1
)

type OptionType uint8

const(
	BackLog		OptionType = 1
	BufferLength	OptionType = 1<<1
	NetWorkType	OptionType = 1<<2
	StreamType	OptionType = 1<<3
)

type ConnectStatus uint8

const(
	NORMAL ConnectStatus = 1
	CLOSED ConnectStatus = 1<<1
	RESETD ConnectStatus = 1<<2
	WARN0  ConnectStatus = 1<<3
	WARN1  ConnectStatus = 1<<4
	WARN2  ConnectStatus = 1<<5
)

type ChannelType uint8

const(
	ChildChannel	ChannelType = 1
	ParentChannel	ChannelType = 1<<1
)

type NetworkInet64 struct{
	Port	uint32
	Address string
}
