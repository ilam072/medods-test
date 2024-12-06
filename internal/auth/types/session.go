package types

import "time"

type Session struct {
	SessionId    string
	UserId       string
	RefreshToken string
	ExpiresAt    time.Time
	Used         bool
}

func (s *Session) IsRefreshTokenExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
