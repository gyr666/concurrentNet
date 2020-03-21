package netpoll

type UnixPoll interface {
	Handle(h func(interface{}))
	Select()
	Stop()
}

type UnixServerPoll interface {
	Init(ServerChannel) UnixServerPoll
	UnixPoll
}
