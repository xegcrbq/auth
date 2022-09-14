package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/models"
)

const (
	UserTable = "users"
)

type CredentialsRepo struct {
	db *sqlx.DB
}

func NewCredentialsRepo(db *sqlx.DB) *CredentialsRepo {
	return &CredentialsRepo{
		db: db,
	}
}
func (sr *CredentialsRepo) ReadCredentialsByUsername(username string) (*models.Credentials, error) {
	var refreshSessions models.Credentials
	sql, _, _ := sq.Select("*").From(UserTable).Where("username=$1").ToSql()
	err := sr.db.Get(&refreshSessions, sql, username)
	return &refreshSessions, err
}

func (sr *CredentialsRepo) SaveCredentials(c *models.Credentials) error {
	sql, _, _ := sq.Insert(UserTable).Columns("username", "password").Values(c.Username, c.Password).ToSql()
	_, err := sr.db.Exec(sql, c.Username, c.Password)
	return err
}

func (sr *CredentialsRepo) DeleteCredentialsByUsername(username string) error {
	sql, _, _ := sq.Delete(UserTable).Where("username = $1").ToSql()
	_, err := sr.db.Exec(sql, username)
	return err
}
