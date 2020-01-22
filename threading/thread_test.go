package threading

import (
	"fmt"
	"testing"
	"time"
)

func TestThreadPoolImpl_Boot(t *testing.T) {
	tp := NewThreadPool(20, 5, time.Second, 10, nil)
	tp.ShutdownNow()
	tp.WaitForStop()
	fmt.Print("Ddd")
}

func TestThreadPoolImpl_Exec(t *testing.T) {
	tp := NewThreadPool(10, 5, time.Second, 10, nil)
	f := tp.Execwpr(func(i ...interface{}) interface{} {
		return i[0].(int) + i[1].(int)
	}, 123, 456)
	t.Logf("result =%d", f.get())
	tp.ShutdownNow()
	fmt.Print("Ddd")
}
