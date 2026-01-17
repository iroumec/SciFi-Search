package structures

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"container/list"
	"scifi-search/app/utils/converters"
	"slices"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// List
func ListContains(l *list.List, value any) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return true
		}
	}
	return false
}

// ------------------------------------------------------------------------------------------------

func AddIfNotExists(l *list.List, value any) {
	if !ListContains(l, value) && converters.ToString(value) != "" {
		l.PushFront(value)
	}
}

// ------------------------------------------------------------------------------------------------

// Array
func Exists[T comparable](a []T, value T) bool {
	return slices.Contains(a, value)
}

// ------------------------------------------------------------------------------------------------
