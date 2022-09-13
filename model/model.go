package model

import "time"

// Credentials структура для данных пользователя
type Credentials struct {
	UserId   int32  `json:"userId" db:"userId"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"userName"`
}

func (c Credentials) IsValid() bool {
	if c.Username != "" && c.Password != "" {
		return true
	}
	return false
}

// RefreshSession структура для данных сессии
type RefreshSession struct {
	Id          int32     `db:"id"`
	UserId      int32     `db:"userid"`
	ReToken     string    `db:"refreshtoken"`
	UserAgent   string    `db:"useragent"`
	Fingerprint string    `db:"fingerprint"`
	Ip          string    `db:"ip"`
	ExpiresIn   int64     `db:"expiresin"`
	CreatedAt   time.Time `db:"createdat"`
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
func (r RefreshSession) Equal(r2 RefreshSession) bool {
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
