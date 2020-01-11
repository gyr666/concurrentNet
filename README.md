# concurrentNet
a network platform to support high concurrent

# GET START
### A simple echo server
```go
server := NewConcurrentNet()
exeCode := server.OnChannelConnection(func(c Channel,p Pipeline){
	p.addLast(func(d Data) Data){
		return d
	})
}).
	setServerScoketChannel(ChannelFactory.Instance.KqueueSocketChannel).
	option(Option.BackLog,1024).
	option(Option.BufferLength,2020).
	option(Option.NetWorkType,NetWorkType.TCP).
	addLisetn(&NetworkInet46{Port:7788}).
	wtype(Wtype.ASYNC).
	sync()
unix.SIGNAL(2,func(signal int)){
	server.stop()
})
fmt.Println("server listen at 0.0.0.0:7788,[::]:7788")
sleep(-1)

```