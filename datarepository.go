package auth

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
)

type DbCredentials struct {
	host     string
	port     int64
	user     string
	password string
	dbname   string
}

func NewDefaultDbCredentials() DbCredentials {
	return DbCredentials{
		host:     "ec2-34-243-101-244.eu-west-1.compute.amazonaws.com",
		port:     5432,
		user:     "hvbofdxjbkkdgq",
		password: "ff9c8195d4fa5205036cb92a384e142c9ca7bfbbc5f7639f038b4925bacdfea9",
		dbname:   "d62omvefcmhpmq",
	}
}
func NewDbCredentials(host string, port int64, user, password, dbname string) DbCredentials {
	return DbCredentials{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

func (dbC DbCredentials) dbDataSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		dbC.host, dbC.port, dbC.user, dbC.password, dbC.dbname)
}

func (dbC DbCredentials) DbCreateTables() error {
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE users(
	  "userId" SERIAL PRIMARY KEY ,
	  "userName" varchar(50) NOT NULL,
	  "password" varchar(50) NOT NULL
	);
	CREATE TABLE refreshSessions(
		"id" SERIAL PRIMARY KEY,
		"userId" integer REFERENCES users("userId") ON DELETE CASCADE,
		"refreshToken" varchar(300) NOT NULL,
		"ua" character varying(200) NOT NULL, /* user-agent */
		"fingerprint" varchar(300) NOT NULL,
		"ip" character varying(15) NOT NULL,
		"expiresIn" bigint NOT NULL,
		"createdAt" timestamp with time zone NOT NULL DEFAULT now()
	);
`)
	if err != nil {
		return err
	}
	return nil
}
func (dbC DbCredentials) addData(data dbData) error {
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		return err
	}
	switch data.(type) {
	case model.Credentials:
		creds := data.(model.Credentials)
		_, err = db.Exec(`
		INSERT INTO users ("userName", "password")
		VALUES ($1, $2);`,
			creds.Username, creds.Password)
	case RefreshSession:
		refreshSession := data.(RefreshSession)
		_, err = db.Exec(`INSERT INTO refreshSessions ("id", "userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn", "createdAt") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
			refreshSession.Id, refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint,
			refreshSession.Ip, refreshSession.ExpiresIn, refreshSession.CreatedAt)
	}
	return err
}

func (dbC DbCredentials) removeData(data dbData) error {
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		return err
	}
	switch data.(type) {
	case model.Credentials:
		creds := data.(model.Credentials)
		_, err = db.Exec(`
		INSERT INTO users ("userName", "password")
		VALUES ($1, $2);`,
			creds.Username, creds.Password)
	case RefreshSession:
		refreshSession := data.(RefreshSession)
		_, err = db.Exec(`INSERT INTO refreshSessions ("id", "userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn", "createdAt") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
			refreshSession.Id, refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint,
			refreshSession.Ip, refreshSession.ExpiresIn, refreshSession.CreatedAt)
	}
	return err
}
