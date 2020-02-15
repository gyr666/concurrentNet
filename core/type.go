package core

type WaitType uint8

const (
	SYNC  WaitType = 1
	ASYNC WaitType = 1 << 1
)

type ConnectStatus uint8

const (
	NORMAL ConnectStatus = 1
	CLOSED ConnectStatus = 1 << 1
	RESET  ConnectStatus = 1 << 2
	WARN0  ConnectStatus = 1 << 3
	WARN1  ConnectStatus = 1 << 4
	WARN2  ConnectStatus = 1 << 5
)

type NetworkType uint8

const (
	TCP NetworkType = 1
	UDP NetworkType = 1 << 1
	RAW NetworkType = 1 << 2
)

type ChannelType uint8

const (
	Child  ChannelType = 1
	Parent ChannelType = 1 << 1
)

type NetworkInet64 struct {
	Port    uint32
	Address string
}