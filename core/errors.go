package core

import "errors"

var(
	errLoopShutdown = errors.New("loop is going to be shutdown")
	errBoot = errors.New("boot error")
	errConnClose = errors.New("this connection is closed")
)