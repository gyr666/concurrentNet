package core

type CallBackEvent interface {
	ReadEventComplete(channelId uint64)
	ReadEventAsyncExecuteComplete(channelId uint64)
	WriteEventComplete(channelId uint64)
	ReadPipelineComplete(channelId uint64)
	WriteEventAsyncExecuteComplete(channelId uint64)
	PreWriteException(channelId uint64, err error)
	PreReadException(channelId uint64, err error)
	OperatorException(channelId uint64, err error)
	PeerReset(channelId uint64)
}
