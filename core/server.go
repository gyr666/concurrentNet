package core

import "fmt"

type Server interface{
	OnChannelConnect(func(c Channel,p Pipeline)) Server
	SetServerScoketChannel(ServerInstance) Server
	Option(OptionType,interface{}) Server
	AddListen(*NetworkInet64) Server
	Wtype(WaitType) Server
	Stop()
	Sync() uint8
	Join()
}


type ServerImpl struct{
	u chan uint8
}
func (s*ServerImpl) Init() Server{
	s.u = make(chan uint8,1)
	return s
}
func (s*ServerImpl) OnChannelConnect(func(c Channel,p Pipeline)) Server{
	return s
}
func (s*ServerImpl) SetServerScoketChannel(ServerInstance) Server{
	return s
}
func (s*ServerImpl) Option(OptionType,interface{}) Server{
	return s
}
func (s*ServerImpl) AddListen(*NetworkInet64) Server{
	return s
}
func (s*ServerImpl) Wtype(WaitType) Server{
	return s
}
func (s*ServerImpl) Stop(){
	s.u <- 1
	fmt.Print("System Stoping\n");
}
func (s*ServerImpl) Join(){
}
func (s*ServerImpl) Sync() uint8 {
	//todo use callback
	fmt.Print("System Booting\n");
	<-s.u
	return 0
}
