package threading

import (
	"fmt"
	"testing"
	"time"
)

func TestThreadpoolImpl_Boot(t *testing.T) {
	tp :=NewThreadPool(2,5,time.Now(),nil)
	tp.Exec(func() {
		for i:=0;i<1 ;i++  {
			fmt.Println("hello world")
		}
	})
	f := tp.Execwr(func() interface{} {
		time.Sleep(6666666)
		k := 8
		for i:=0;i<9292 ;i++  {
			k +=i
		}
		return k
	})
	fmt.Println( f.get())
	tp.Shutdown()
	fmt.Print("Ddd")
}