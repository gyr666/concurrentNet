package core

type Server interface{
	OnChannelConnection(func(c Channel,p Pipeline)) Server
	SetServerScoketChannel(ServerInstance) Server
	Option(OptionType,interface{}) Server
	AddLestion(NetworkInet64) Server
	Wtype(WaitType) Server
	Join()
}