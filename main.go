package main

import (
	"golang.org/x/sys/unix"
	"fmt"
	"gunplan.top/concurrentNet/core"
)

func main() {
	sig :=make(chan os.Signal
	server  := core.NewConcurrentNet()
	server.OnChannelConnect(func(c core.Channel,p core.Pipeline){
		p.AddLast(func(d core.Data) core.Data{
			return d
		})
	}).
	SetServerScoketChannel(ChannelFactory.Instance.KqueueSocketChannel).
	Option(&core.BackLog{},1024).
	Option(&core.BufferLength{},2020).
	Option(&core.NetWorkType{},NetWorkType.TCP).
	AddListen(&core.NetworkInet64{Port:7788}).
	Wtype(core.ASYNC).
	Sync()
	unix.Signal(2,func(signal int){
		server.Stop()
	})
	fmt.Println("server listen at 0.0.0.0:7788,[::]:7788")
	server.Join()
}
