package disruptor

import (
	"fmt"
	"testing"
	"time"
)

func TestRing(t *testing.T) {
	d := &disruptorImpl{}
	d.Init(20)
	for i := 0; i < 5; i++ {
		d.Push(i)
	}

	for i := 0; i < 5; i++ {
		fmt.Println(d.Poll())
	}
}

func TestRingMixedThread(t *testing.T) {
	d := &disruptorImpl{}
	d.Init(20)
	for i := 0; i < 5; i++ {
		d.Push(i)
	}
	time.Sleep(1000)
	for j := 0; j < 2; j++ {
		go func() {
			for i := 0; i < 3; i++ {
				fmt.Println(d.Poll())
			}
		}()
	}
	time.Sleep(1000)

}
