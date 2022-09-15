package repositories

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/models"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

const (
	SessionTable = "refreshsessions"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}
func (sr *SessionRepo) ReadSessionByRefreshToken(cmd *models.QueryReadSessionByRefreshToken) (*models.Session, error) {
	var refreshSessions models.Session
	sql, _, _ := sq.Select("*").From(SessionTable).Where("refreshtoken=$1").ToSql()
	err := sr.db.Get(&refreshSessions, sql, cmd.RefreshToken)
	return &refreshSessions, err
}

func (sr *SessionRepo) SaveSession(cmd *models.CommandCreateSession) error {
	rs := cmd.Session
	sql, _, _ := sq.Insert(SessionTable).Columns("userid", "refreshtoken", "useragent", "fingerprint", "ip", "expiresin", "createdat").Values(rs.UserId, rs.ReToken, rs.UserAgent, rs.Fingerprint, rs.Ip, rs.ExpiresIn, rs.CreatedAt).ToSql()
	_, err := sr.db.Exec(sql, rs.UserId, rs.ReToken, rs.UserAgent, rs.Fingerprint, rs.Ip, rs.ExpiresIn, rs.CreatedAt)
	return err
}

func (sr *SessionRepo) DeleteSessionByRefreshToken(cmd *models.CommandDeleteSessionByRefreshToken) error {
	sql, _, _ := sq.Delete(SessionTable).Where("refreshtoken = $1").ToSql()
	_, err := sr.db.Exec(sql, cmd.RefreshToken)
	return err
}
func (sr *SessionRepo) GetCommand(comand interface{}) {

}
