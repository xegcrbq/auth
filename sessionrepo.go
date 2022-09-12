package auth

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

func (repo *SessionRepo) GetExistingData(reToken string) (RefreshSession, error) {
	var reSession RefreshSession
	repo.db.Get(&reSession, `SELECT * FROM refreshsessions WHERE "refreshToken" = $1;`, reToken)
	if reSession.ExpiresIn != 0 {
		return reSession, nil
	}
	return reSession, errors.New("not found")
}

func (repo *SessionRepo) RemoveData(reToken string) (RefreshSession, error) {
	reSession, err := repo.GetExistingData(reToken)
	if err != nil {
		return reSession, err
	}
	res, err := repo.db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, reToken)
	rAffected, err := res.RowsAffected()
	if rAffected > 0 && err == nil {
		return reSession, err
	}
	return reSession, errors.New("delete error")
}
func (repo *SessionRepo) AddData(session RefreshSession) (RefreshSession, error) {
	_, err := repo.db.Exec(`INSERT INTO refreshSessions ("id", "userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn", "createdAt") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		session.Id, session.UserId, session.ReToken, session.UserAgent, session.Fingerprint,
		session.Ip, session.ExpiresIn, session.CreatedAt)
	return session, err
}
