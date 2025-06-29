package models

import "time"

type SessionData struct {
	Token     string
	ExpiresAt time.Time
}
