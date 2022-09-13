package model

import "time"

type RefreshSession struct {
	Id          int32     `db:"id"`
	UserId      int32     `db:"userId"`
	ReToken     string    `db:"refreshToken"`
	UserAgent   string    `db:"ua"`
	Fingerprint string    `db:"fingerprint"`
	Ip          string    `db:"ip"`
	ExpiresIn   int64     `db:"expiresIn"`
	CreatedAt   time.Time `db:"createdAt"`
}

func (r RefreshSession) IsValid() bool {
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
