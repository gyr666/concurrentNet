package netpoll

import "golang.org/x/sys/unix"

//todo

type serverUnixPollImpl struct {
	fd     int
	handle func(interface{})
}

type ServerChannel struct {
	fd int32
}

func (s *ServerChannel) Id32() int32 {
	return s.fd
}

func (s *ServerChannel) Id() int {
	return int(s.fd)
}

func (s *serverUnixPollImpl) Init(server ServerChannel) UnixServerPoll {
	var err error

	if s.fd, err = unix.EpollCreate1(0); err != nil {
		panic(err)
	}

	e := new(unix.EpollEvent)
	e.Events = unix.EPOLLIN
	e.Fd = server.Id32()

	if unix.EpollCtl(s.fd, unix.EPOLL_CTL_ADD, server.Id(), e) != nil {
		panic(err)
	}
	return s
}

func (s *serverUnixPollImpl) Handle(h func(interface{})) {
	s.handle = h
}

func (s *serverUnixPollImpl) Select() {

}

func (s *serverUnixPollImpl) Stop() {

}
