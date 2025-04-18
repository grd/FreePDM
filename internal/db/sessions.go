package db

import "time"

type Session struct {
	Username   string
	Expiration time.Time
}

var Sessions = map[string]Session{}
