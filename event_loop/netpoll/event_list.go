package netpoll

import "golang.org/x/sys/unix"

type eventList struct {
	events []unix.EpollEvent
	size int
}

func NewEventList(size int)*eventList{
	return  &eventList{size:size,events:make([]unix.EpollEvent,size)}
}

func (el *eventList)increase(){
	el.size<<=1
	el.events = make([]unix.EpollEvent,el.size)
}