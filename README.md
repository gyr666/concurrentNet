# ConcurrentNet
a network platform to support high concurrent

# GET START
### A simple echo server
```go
server  := NewConcurrentNet()
exeCode := server.OnChannelConnect(func(c Channel,p Pipeline){
	p.addLast(func(d Data) Data){
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
unix.SIGNAL(2,func(signal int)){
	server.Stop()
})
fmt.Println("server listen at 0.0.0.0:7788,[::]:7788")
server.Join()
```

### Then Test It
```bash
telnet [::]:1 7788
telent> Hello
telnet> Hello
```

# PACKAGE INTURDUCE
## buffer
package `buffer` implement the `buffer alloctor` by `bitmap` and `skiplist`.
## pipeline
package `pipeline` implement pipe that when channel create, data I/O, exception, channel close.
## threading
package `threading` implement thread pool by stand threadpool and steal task threadpool

