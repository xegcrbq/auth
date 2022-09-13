package sqlstore

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
)

const (
	UserTable    = "users"
	SessionTable = "refreshsessions"
)

func DropTables(db *sqlx.DB) {
	db.Exec(`DROP TABLE ` + UserTable)
	db.Exec(`DROP TABLE ` + SessionTable)
}

func CreateTables(db *sqlx.DB) {
	db.Exec(`CREATE TABLE ` + UserTable + `(
	userid SERIAL PRIMARY KEY ,
	username varchar(50) NOT NULL,
	password varchar(50) NOT NULL
	)`)
	db.Exec(`CREATE TABLE ` + SessionTable + `(
    id SERIAL PRIMARY KEY,
    userid integer REFERENCES users(userid) ON DELETE CASCADE,
    refreshtoken varchar(300) NOT NULL,
    useragent character varying(200) NOT NULL,
    fingerprint varchar(300) NOT NULL,
    ip character varying(15) NOT NULL,
    expiresin bigint NOT NULL,
    createdAt timestamp with time zone NOT NULL DEFAULT now()
	)`)
}

type SqlxDatabase interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func ReadByRefreshToken(ctx context.Context, db SqlxDatabase, refreshToken string) ([]*model.RefreshSession, error) {
	var refreshSessions []*model.RefreshSession
	sql := `SELECT * FROM ` + SessionTable + ` WHERE refreshtoken=$1`
	err := db.SelectContext(ctx, &refreshSessions, sql, refreshToken)
	return refreshSessions, err
}

func Save(ctx context.Context, db SqlxDatabase, rs *model.RefreshSession) error {
	sql := `INSERT INTO ` + SessionTable + ` (userid, refreshtoken, useragent, fingerprint, ip, expiresin, createdat) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.ExecContext(ctx, sql, rs.UserId, rs.ReToken, rs.UserAgent, rs.Fingerprint, rs.Ip, rs.ExpiresIn, rs.CreatedAt)
	return err
}
func DeleteByRefreshToken(ctx context.Context, db SqlxDatabase, refreshToken string) error {
	sql := `DELETE FROM ` + SessionTable + ` WHERE refreshtoken = $1`
	_, err := db.ExecContext(ctx, sql, refreshToken)
	return err
}
