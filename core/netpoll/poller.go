package netpoll

import (
	"golang.org/x/sys/unix"
	"gunplan.top/concurrentNet/util"
	"log"
	"unsafe"
)

type Poller struct {
	pfd    int
	wfd    int
	wfdBuf []byte
	queue  *util.Queue
}

func NewPoller() (*Poller, error) {
	pfd, err0 := unix.EpollCreate1(0)
	if err0 != nil {
		return nil, err0
	}
	wfd, err := unix.Eventfd(0, 0)
	if err != nil {
		return nil, err
	}
	p := new(Poller)
	p.pfd = pfd
	p.wfd = wfd
	p.wfdBuf = make([]byte, 8)
	p.queue = util.NewQueue()
	if err = p.AddReader(p.wfd); err != nil {
		return nil, err
	}
	return p, nil
}

func(p *Poller)Close(){
	_ =unix.Close(p.wfd)
	_= unix.Close(p.pfd)
}

func (p *Poller) Polling(callback func() error) error {
	var wakeup bool
	el := newEventList(128)
	for {
		n, err0 := unix.EpollWait(p.pfd, el.events, -1)
		if err0 != nil && err0 == unix.EINTR {
			log.Println(err0)
			continue
		}
		for i := 0; i < n; i++ {
			if fd := int(el.events[i].Fd); fd != p.wfd {
				if err := callback(); err != nil {
					log.Println(err)
				}
			} else {
				wakeup = true
				_, _ = unix.Read(p.wfd, p.wfdBuf)
			}
			if wakeup {
				wakeup = false
				if err := p.queue.ForEach(); err != nil {
					return err
				}
			}
		}
		if n == el.size {
			el.increase()
		}
	}
}

var (
	u uint64 = 1
	b []byte = (*(*[8]byte)(unsafe.Pointer(&u)))[:]
)

func (p *Poller) Trigger(callback func() error) error {
	if num := p.queue.Push(callback); num == 1 {
		_, err := unix.Write(p.wfd, b)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	readEvent      = unix.EPOLLPRI | unix.EPOLLIN
	writeEvent     = unix.EPOLLOUT
	readWriteEvent = readEvent | writeEvent
)

func (p *Poller) AddReader(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvent})
}

func (p *Poller) AddWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: writeEvent})
}

func (p *Poller) AddReadWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readWriteEvent})
}

func (p *Poller) ModReader(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvent})
}

func (p *Poller) ModWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: writeEvent})
}

func (p *Poller) ModReadWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readWriteEvent})
}
