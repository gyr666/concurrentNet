package main

import (
	"fmt"
	"gunplan.top/concurrentNet/core"
)

func main() {
	server  := core.NewConcurrentNet()
	server.OnChannelConnect(func(c core.Channel,p core.Pipeline){
		p.addLast(func(d core.Data) core.Data{
			return d
		})
	}).
	SetServerScoketChannel(ChannelFactory.Instance.KqueueSocketChannel).
	Option(Option.BackLog,1024).
	Option(Option.BufferLength,2020).
	Option(Option.NetWorkType,NetWorkType.TCP).
	AddListen(&NetworkInet64{Port:7788}).
	Wtype(Wtype.ASYNC).
	Sync()
	unix.Signal(2,func(signal int){
		server.Stop()
	})
	fmt.Println("server listen at 0.0.0.0:7788,[::]:7788")
	server.Join()
}
