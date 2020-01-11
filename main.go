package main

import (
	"fmt"
	"gunplan.top/concurrentNet/event_loop"
)

func main() {
	event_loop.NewBaseEventLoop().StartLoop()
	fmt.Print("listen 7887")
}
