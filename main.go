package main

import (
	"os"
	"os/signal"
	"syscall"

	"gunplan.top/concurrentNet/config"
	"gunplan.top/concurrentNet/core"
)

func main() {
	c := config.GetConfig{C: config.LineDecoder, F: config.DefaultFetcher}
	server := core.NewConcurrentNet()
	server.OnChannelConnect(func(c core.Channel, p core.Pipeline) {
		p.AddLast(func(d core.Data) core.Data {
			return d
		})
	}).
		Option(&c).
		WaitType(config.ASYNC)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sc
		server.Stop()
	}()
	_ = server.Sync()
	server.Join()
}
