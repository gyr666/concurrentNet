package netpoll

import "golang.org/x/sys/unix"

func NewServerChannel(port, backlog int) ServerChannel {
	serverFd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	err = unix.Bind(serverFd, &unix.SockaddrInet4{Port: port})
	if err != nil {
		panic(err)
	}

	err = unix.Listen(serverFd, backlog)
	if err != nil {
		panic(err)
	}
	return ServerChannel{fd: int32(serverFd)}
}
