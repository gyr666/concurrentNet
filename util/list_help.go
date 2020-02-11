package util

import "container/list"

func Traverse(list *list.List, f func(element *list.Element) bool) {
	for item := list.Front(); nil != item; item = item.Next() {
		if f(item) {
			list.Remove(item)
			break
		}
	}
}
