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

func (r *RefreshSession) SetId(newId int32) {
	r.Id = newId
}
func (r *RefreshSession) SetUserId(newUserId int32) {
	r.UserId = newUserId
}
func (r *RefreshSession) SetRefreshToken(newToken string) {
	r.ReToken = newToken
}
func (r *RefreshSession) SetUserAgent(newUserAgent string) {
	r.UserAgent = newUserAgent
}
func (r *RefreshSession) SetFingerprint(newFingerprint string) {
	r.Fingerprint = newFingerprint
}
func (r *RefreshSession) SetIp(newIp string) {
	r.Ip = newIp
}
func (r *RefreshSession) SetExpiresIn(newExpiresIn int64) {
	r.ExpiresIn = newExpiresIn
}
func (r *RefreshSession) SetCreatedAt(newCreatedAt time.Time) {
	r.CreatedAt = newCreatedAt
}
func (r RefreshSession) GetId() int32 {
	return r.Id
}
func (r RefreshSession) GetUserId() int32 {
	return r.UserId
}
func (r RefreshSession) GetRefreshToken() string {
	return r.ReToken
}
func (r RefreshSession) GetUserAgent() string {
	return r.UserAgent
}
func (r RefreshSession) GetFingerprint() string {
	return r.Fingerprint
}
func (r RefreshSession) GetIp() string {
	return r.Ip
}
func (r RefreshSession) GetExpiresIn() int64 {
	return r.ExpiresIn
}
func (r RefreshSession) GetCreatedAt() time.Time {
	return r.CreatedAt
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
