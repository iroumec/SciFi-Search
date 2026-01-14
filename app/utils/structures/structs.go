package structures

import (
	"time"
)

type HistoricSearchView struct {
	Query string
	Date  time.Time
}

type User struct {
	Name            string
	Surname         string
	AvatarURLString string
	AvatarURLValid  bool
	Email           string
}
