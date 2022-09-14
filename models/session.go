package models

import "time"

// Session структура для данных сессии
type Session struct {
	Id          int32     `db:"id"`
	UserId      int32     `db:"userid"`
	ReToken     string    `db:"refreshtoken"`
	UserAgent   string    `db:"useragent"`
	Fingerprint string    `db:"fingerprint"`
	Ip          string    `db:"ip"`
	ExpiresIn   int64     `db:"expiresin"`
	CreatedAt   time.Time `db:"createdat"`
}

func (r Session) IsValid() bool {
	if r.UserId == 0 {
		return false
	}
	if r.ReToken == "" {
		return false
	}
	if r.UserAgent == "" {
		return false
	}
	if r.Fingerprint == "" {
		return false
	}
	if r.Ip == "" {
		return false
	}
	if r.ExpiresIn == 0 {
		return false
	}
	if r.CreatedAt == time.UnixMicro(0) {
		return false
	}
	return true
}
func (r Session) Equal(r2 Session) bool {
	if r.UserId != r2.UserId {
		return false
	}
	if r.ReToken != r2.ReToken {
		return false
	}
	if r.UserAgent != r2.UserAgent {
		return false
	}
	if r.Fingerprint != r2.Fingerprint {
		return false
	}
	if r.Ip != r2.Ip {
		return false
	}
	if r.ExpiresIn != r2.ExpiresIn {
		return false
	}
	if r.CreatedAt.Round(time.Second).Unix() != r2.CreatedAt.Round(time.Second).Unix() {
		return false
	}
	return true
}

type RefreshSessionRepository interface {
	SaveSession(rs *Session) error
	ReadSessionByRefreshToken(refreshToken string) (*Session, error)
	DeleteSessionByRefreshToken(refreshToken string) error
}
