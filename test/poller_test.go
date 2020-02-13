package test

import (
	"fmt"
	"gunplan.top/concurrentNet/event_loop/netpoll"
	"sync"
	"testing"
)

func TestPoller(t *testing.T){
	p,err:=netpoll.NewPoller()
	if err != nil{
		t.Log("poller init error!")
	}
	wg:=sync.WaitGroup{}
	wg.Add(1)

	go p.Polling(func(int)error{
		return nil
	})

	p.Trigger(func()error{
		t.Log("trigger!")
		return nil
	})
	p.Trigger(func()error{
		wg.Done()
		return fmt.Errorf("")
	})
	wg.Wait()
}