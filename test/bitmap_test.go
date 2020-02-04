package test

import (
	"fmt"
	"testing"
)
import "gunplan.top/concurrentNet/util"

import "C"

func TestBitMap(t *testing.T){
	 var bm C.BitMap = C.BitMapInit(29);
	 bm.MaktItAs(bm.data,1)
	 fmt.Print(AcquirePosition(bm.data))
}
