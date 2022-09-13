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
func (r1 *RefreshSession) Equal(r2 RefreshSession) bool {
	if r1.UserId != r2.UserId {
		return false
	}
	if r1.ReToken != r2.ReToken {
		return false
	}
	if r1.UserAgent != r2.UserAgent {
		return false
	}
	if r1.Fingerprint != r2.Fingerprint {
		return false
	}
	if r1.Ip != r2.Ip {
		return false
	}
	if r1.ExpiresIn != r2.ExpiresIn {
		return false
	}
	if r1.CreatedAt.Round(time.Second).Unix() != r2.CreatedAt.Round(time.Second).Unix() {
		return false
	}
	return true
}
