package core

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
