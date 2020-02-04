package core

type Channel interface{
	Address() NetworkInet64
	Status() ConnectStatus
	Write(Data)
	Read() Data
	Close() error
	Reset() error
	Type()  ChannelType
	parent() Channel
}
