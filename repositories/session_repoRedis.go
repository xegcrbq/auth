package repositories

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/xegcrbq/auth/models"
	"time"
)

type SessionRepoRedis struct {
	db *redis.Client
}

func NewSessionRepoRedis(db *redis.Client) *SessionRepoRedis {
	return &SessionRepoRedis{
		db: db,
	}
}
func (sr *SessionRepoRedis) ReadSessionByRefreshToken(cmd *models.QueryReadSessionByRefreshToken) (*models.Session, error) {
	var refreshSessions models.Session
	val, err := sr.db.Get(cmd.RefreshToken).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(val), &refreshSessions)
	return &refreshSessions, err
}

func (sr *SessionRepoRedis) SaveSession(cmd *models.CommandCreateSession) error {
	rs := cmd.Session
	json, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	err = sr.db.Set(rs.ReToken, json, time.Until(time.Unix(rs.ExpiresIn, 0))).Err()
	return err
}

func (sr *SessionRepoRedis) DeleteSessionByRefreshToken(cmd *models.CommandDeleteSessionByRefreshToken) error {
	return sr.db.Del(cmd.RefreshToken).Err()
}
