# ConcurrentNet
a network platform to support high concurrent

# GET START
### A simple echo server
```go
server  := NewConcurrentNet()
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
server.join()
```

### Then Test It
```bash
telnet [::]:1 7788
telent> Hello
tlenet> Hello
```

# PACKAGE INTURDUCE
## buffer
package `buffer` implement the `buffer alloctor` by `bitmap` and `skiplist`.
## pipeline
package `pipeline` implement pipe that when channel create, data I/O, exception, channel close.
## threading
package `threading` implement thread pool by stand threadpool and steal task threadpool

