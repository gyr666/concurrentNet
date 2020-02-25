package core

import "strings"

type ConnectStatus uint8

const (
	NORMAL ConnectStatus = 1
	CLOSED ConnectStatus = 1 << 1
	RESET  ConnectStatus = 1 << 2
	WARN0  ConnectStatus = 1 << 3
	WARN1  ConnectStatus = 1 << 4
	WARN2  ConnectStatus = 1 << 5
)

type ServerStatus uint8

const(
	NONE ServerStatus = iota
	BOOTING
	RUNNING
	STOPPING
	STOPPED
)

type NetworkInet64 struct {
	network   string
	Address string
}

func parseAddress(addr string)(network,address string){
	network="tcp"
	address=addr
	if strings.Contains(addr,"://"){
		parts := strings.Split(addr,"://")
		network=parts[0]
		address=parts[1]
	}
	return
}
