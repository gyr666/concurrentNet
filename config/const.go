package config

import "gunplan.top/concurrentNet/buffer"

type NetworkType uint8

const (
	TCP NetworkType = 1
	UDP NetworkType = 1 << 1
	RAW NetworkType = 1 << 2
	RTP NetworkType = 1 << 3
	RTCP NetworkType = 1 << 4
)

type WaitType uint8

const (
	SYNC  WaitType = 1
	ASYNC WaitType = 1 << 1
)

type Config struct {
	Backlog			uint8
	BufferSize 		uint16
	AllocType	 	buffer.Allocator
	NetworkType 	NetworkType
	LogFile			string
	Acl				uint8
	WaitType		WaitType
	MaxConnection 	uint64
	TimerTick		uint8
}