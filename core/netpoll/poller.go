package netpoll

import (
	"log"
	"unsafe"

	"golang.org/x/sys/unix"

	"gunplan.top/concurrentNet/util"
)

type EventCallback func(int, uint32) error

type Polling struct {
	c      EventCallback
	pfd    int
	wfd    int
	wfdBuf []byte
	queue  *util.Queue
}

func NewPolling(c EventCallback) *Polling {
	pfd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	wfd, err := unix.Eventfd(0, 0)
	if err != nil {
		panic(err)
	}
	p := new(Polling)
	p.pfd = pfd
	p.wfd = wfd
	p.wfdBuf = make([]byte, 8)
	p.c = c
	p.queue = util.NewQueue()
	if err = p.AddReader(p.wfd); err != nil {
		panic(err)
	}
	return p
}

func (p *Polling) Close() {
	_ = unix.Close(p.wfd)
	_ = unix.Close(p.pfd)
}

func (p *Polling) Polling() error {
	var wakeup bool
	el := newEventList(InitEvents)
	for {
		n, err0 := unix.EpollWait(p.pfd, el.events, -1)
		if err0 != nil && err0 == unix.EINTR {
			log.Println(err0)
			continue
		}
		for i := 0; i < n; i++ {
			if fd, events := int(el.events[i].Fd), el.events[i].Events; fd != p.wfd {
				if err := p.c(fd, events); err != nil {
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
	b        = (*(*[8]byte)(unsafe.Pointer(&u)))[:]
)

func (p *Polling) Trigger(callback func() error) error {
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

func (p *Polling) AddReader(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvent})
}

func (p *Polling) AddWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: writeEvent})
}

func (p *Polling) AddReadWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readWriteEvent})
}

func (p *Polling) ModReader(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvent})
}

func (p *Polling) ModWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: writeEvent})
}

func (p *Polling) ModReadWriter(fd int) error {
	return unix.EpollCtl(p.pfd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readWriteEvent})
}
