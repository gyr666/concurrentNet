package core

type CallBackEvent interface {
	ChannelReadEventComplete(channelId uint64)
	ChannelReadEventAsyncExecuteComplete(channelId uint64)
	ChannelWriteEventComplete(channelId uint64)
	ChannelReadPipelineComplete(channelId uint64)
	ChannelWriteEventAsyncExecuteComplete(channelId uint64)
	ChannelPreWriteException(channelId uint64, err error)
	ChannelPreReadException(channelId uint64, err error)
	ChannelOperatorException(channelId uint64, err error)
	ChannelPeerReset(channelId uint64)
}
