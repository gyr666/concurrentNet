package test

import (
	"testing"

	"gunplan.top/concurrentNet/core"
)

func TestConnCache(t *testing.T) {
	a := core.NewChannelCache()
	for i := 0; i < 2000000; i++ {
		a.Acquire()
	}
	//for i := 0; i < 20000000; i++ {
	//	core.Factory.NewChannel()
	//}
}

// a allocator and 1 list 11.45s
// a allocator and 20 list 0.80s
// a allocator and 30 list 0.64s
// a allocator and 40 list 0.56s
// no allocator is 0s
