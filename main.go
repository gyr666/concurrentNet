package main

import (
	"os"
	"syscall"
	"os/signal"
	"fmt"
	"gunplan.top/concurrentNet/core"
)

func main() {
	server  := core.NewConcurrentNet()
	server.OnChannelConnect(func(c core.Channel,p core.Pipeline){
		p.AddLast(func(d core.Data) core.Data{
			return d
		})
	}).
	SetServerScoketChannel(core.Factory.NewParentChannelInstance()).
	Option(&core.BackLog{},1024).
	Option(&core.BufferLength{},2020).
	Option(&core.NetWorkType{},core.TCP).
	AddListen(&core.NetworkInet64{Port:7788}).
	Wtype(core.ASYNC)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func(){
		<-sc
		server.Stop()
	}()
	server.Sync()
	fmt.Println("server listen at 0.0.0.0:7788,[::]:7788")
	server.Join()
}
