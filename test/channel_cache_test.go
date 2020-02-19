package test

import (
	"gunplan.top/concurrentNet/core"
	"testing"
)

func TestChannelCache(t *testing.T) {
	for i := 0; i < 1000; i++ {
		core.NewChannelCache().Acquire(core.Child).Release()
	}
}
