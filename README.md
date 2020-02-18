# ConcurrentNet
a network platform to support high concurrent 
## Information
Developer   : [**@gyr666**](https://github.com/gyr666) [@H00001(frank)](https://github.com/H00001)  
Language    : Golang 1.13.5, gnuc17  
Environment : Gnu Linux, OS x. not Windows  
Version     : 0.0.0.1  
License     : GUN GENERAL PUBLIC LICENSE 3.0  
Website     : ?

# GET START
### Test
`./config && make && sudo make install`
### A simple echo server
```go
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

