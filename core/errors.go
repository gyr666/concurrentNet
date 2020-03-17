package core

import "errors"

var (
	errLoopShutdown = errors.New("loop is going to be shutdown")
	errBoot         = errors.New("boot error")
	errChannelClose = errors.New("this channel is closed")
)
