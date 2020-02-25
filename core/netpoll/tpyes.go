package netpoll

import "golang.org/x/sys/unix"

const (
	InitEvents  = 128
	ErrEvents   = unix.EPOLLERR | unix.EPOLLHUP | unix.EPOLLRDHUP
	ReadEvents  = ErrEvents | unix.EPOLLIN | unix.EPOLLPRI
	WriteEvents = ErrEvents | unix.EPOLLOUT
)
