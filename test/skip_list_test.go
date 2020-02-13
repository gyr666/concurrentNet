package test

import (
	"gunplan.top/concurrentNet/util"

	"testing"
)

func TestSkipList(t *testing.T) {
	l := util.NewSkipList()
	l.Insert(2, 7)
	l.Insert(5, 9)
	l.Insert(6, 9999)
	d, _ := l.Search(4)
	Equal(d, 5, "divide")

}
