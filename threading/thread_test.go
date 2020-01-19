package threading

import (
	"fmt"
	"testing"
	"time"
)

func TestThreadpoolImpl_Boot(t *testing.T) {
	tp := NewThreadPool(2, 5, time.Second, 10, nil)
	tp.Exec(func() {
		for i := 0; i < 1; i++ {
			fmt.Println("hello world")
		}
	})
	f := tp.Execwr(func() interface{} {
		time.Sleep(6666666)
		k := 8
		for i := 0; i < 9292; i++ {
			k += i
		}
		return k
	})
	fmt.Println(f.get())
	tp.ShutdownNow()
	tp.WaitForStop()
	fmt.Print("Ddd")
}
