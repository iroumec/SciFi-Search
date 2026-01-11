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

//TODO: estos HSV y User se podrían reemplazar con map[string]any y nos ahorramos este archivo
