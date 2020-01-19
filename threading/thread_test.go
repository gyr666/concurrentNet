package threading

import (
	"fmt"
	"testing"
	"time"
)

func TestThreadpoolImpl_Boot(t *testing.T) {
	tp := NewThreadPool(1, 5, time.Second, 10, nil)
	tp.ShutdownNow()
	tp.WaitForStop()
	fmt.Print("Ddd")
}
