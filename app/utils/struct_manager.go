package utils

import "container/list"

// List
func ListContains(l *list.List, value any) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return true
		}
	}
	return false
}

func AddIfNotExists(l *list.List, value any) {
	if !ListContains(l, value) && ToString(value) != "" {
		l.PushFront(value)
	}
}

// Array
func Exists[T comparable](a []T, value T) bool {
	for _, v := range a {
		if v == value {
			return true
		}
	}
	return false
}
