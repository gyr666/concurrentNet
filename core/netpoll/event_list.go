package netpoll

import "golang.org/x/sys/unix"

type eventList struct {
	events []unix.EpollEvent
	size   int
}

func newEventList(size int) *eventList {
	return &eventList{events: make([]unix.EpollEvent, size), size: size}
}

func (el *eventList) increase() {
	el.size <<= 1
	el.events = make([]unix.EpollEvent, el.size)
}
