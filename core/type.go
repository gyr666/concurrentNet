package core

type ConnectStatus uint8

const (
	NORMAL ConnectStatus = 1
	CLOSED ConnectStatus = 1 << 1
	RESET  ConnectStatus = 1 << 2
	WARN0  ConnectStatus = 1 << 3
	WARN1  ConnectStatus = 1 << 4
	WARN2  ConnectStatus = 1 << 5
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
