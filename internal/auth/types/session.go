package types

import "time"

type Session struct {
	SessionId    string
	UserId       int
	RefreshToken string
	ExpiresAt    time.Time
}
