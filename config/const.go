package config

import (
	"gunplan.top/concurrentNet/buffer"
	"strconv"
	"strings"
)

const TypeName = "_type"
const Number = "number"
const String = "string"
const Map = "map_func"

type NetworkType uint8

const (
	TCP  NetworkType = 1
	UDP  NetworkType = 1 << 1
	RAW  NetworkType = 1 << 2
	RTP  NetworkType = 1 << 3
	RTCP NetworkType = 1 << 4
)

type WaitType uint8

const (
	SYNC  WaitType = 1
	ASYNC WaitType = 1 << 1
)

type Config struct {
	Backlog       uint8            `line:"backlog" _type:"number"`
	BufferSize    uint16           `line:"buffer_size" _type:"number"`
	AllocType     buffer.Allocator `line:"allocator"`
	NetworkType   NetworkType      `line:"network_type" _type:"enum"`
	LogFile       string           `line:"log_file" _type:"string"`
	Acl           int64            `line:"acl" _type:"number"`
	Listen        NetworkInet64    `line:"listen" _type:"map_func"`
	WaitType      WaitType
	MaxConnection uint64 `line:"max_connection",type:"number"`
	TimerTick     uint8  `line:"timer_tick",type:"number"`
}

func (c *Config) SetListen(in string) {
	c.Listen = *parseAddress(in)
}

type NetworkInet64 struct {
	network string
	Address string
	port    int
}

func parseAddress(addr string) *NetworkInet64 {
	n := new(NetworkInet64)
	n.network = "tcp"
	if strings.Contains(addr, ":") {
		parts := strings.Split(addr, ":")
		port := parts[len(parts)-1]
		n.Address = addr[0 : len(addr)-len(port)]
		n.port, _ = strconv.Atoi(port)
	}
	return n
}
