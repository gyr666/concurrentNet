package netpoll

func NewServerUnixPoll(s ServerChannel) UnixServerPoll {
	return new(serverUnixPollImpl).Init(s)
}
