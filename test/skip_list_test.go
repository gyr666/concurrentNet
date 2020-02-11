package test

import (
	"fmt"
	"gunplan.top/concurrentNet/util"

	"testing"
)

func TestSkipList(t *testing.T) {
	l := util.NewSkipList()
	l.Insert(2, 7)
	l.Insert(5, 9)
	l.Insert(6, 9999)
	fmt.Print(l.Search(4))

}
