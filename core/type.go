package core

import "strings"

type ChannelStatus uint8

const (
	NORMAL ChannelStatus = 1
	CLOSED ChannelStatus = 1 << 1
	RESET  ChannelStatus = 1 << 2
	WARN0  ChannelStatus = 1 << 3
	WARN1  ChannelStatus = 1 << 4
	WARN2  ChannelStatus = 1 << 5
)

type ServerStatus uint8

const (
	NONE ServerStatus = iota
	BOOTING
	RUNNING
	STOPPING
	STOPPED
)

type NetworkInet64 struct {
	network string
	Address string
}

func parseAddress(addr string) (network, address string) {
	network = "tcp"
	address = addr
	if strings.Contains(addr, "://") {
		parts := strings.Split(addr, "://")
		network = parts[0]
		address = parts[1]
	}
	return
}
