package repositories

import "github.com/xegcrbq/auth/models"

type SessionRepo interface {
	ReadSessionByRefreshToken(cmd *models.QueryReadSessionByRefreshToken) (*models.Session, error)
	SaveSession(cmd *models.CommandCreateSession) error
	DeleteSessionByRefreshToken(cmd *models.CommandDeleteSessionByRefreshToken) error
}
