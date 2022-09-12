package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

// ключ для создания подписи
var jwtKey = []byte("irjgdkngfdkjdkjlbvnjkd")

//var dbDataSource = fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
//	"postgres", 1, "postgres")

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

func (dbC DbCredentials) dbCreateTables() error {
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
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
func (dbC DbCredentials) dbCreateUsers() error {
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	INSERT INTO users ("userName", "password")
	VALUES ('admin', 'admin'),
	('user', 'password');
`)
	if err != nil {
		return err
	}
	return nil
}

type Data interface {
	Valid() error
}

// Credentials структура для парсинга данных из json и бд
type Credentials struct {
	UserId   string `json:"userId" db:"userId"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"userName"`
}

type RefreshSession struct {
	Id          string    `db:"id"`
	UserId      string    `db:"userId"`
	ReToken     string    `db:"refreshToken"`
	UserAgent   string    `db:"ua"`
	Fingerprint string    `db:"fingerprint"`
	Ip          string    `db:"ip"`
	ExpiresIn   int64     `db:"expiresIn"`
	CreatedAt   time.Time `db:"createdAt"`
}

func (c Claims) Valid() error {
	if c == (Claims{}) {
		return errors.New("empty struct")
	} else {
		return nil
	}
}

// Claims структура для генерации из неё jwt токена
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (c Fp) Valid() error {
	if c == (Fp{}) {
		return errors.New("empty struct")
	} else {
		return nil
	}
}

// Fp структура для генерации из неё jwt токена для отпечатка
type Fp struct {
	Fingerprint string `json:"fingerprint"`
	jwt.StandardClaims
}

func CreateJWT(name string, expirationTime time.Time, httpOnly bool, rawData Data, shortPath bool) (*http.Cookie, error) {

	data := jwt.NewWithClaims(jwt.SigningMethodHS256, rawData)
	dataString, err := data.SignedString(jwtKey)
	if err != nil {
		return &http.Cookie{}, errors.New("data.SignedString(jwtKey) error")
	}
	//выставляем параметр HttpOnly, чтобы получать доступ к этому токену только на странице авторизации
	if shortPath {
		return &http.Cookie{
			Name:     name,
			Value:    dataString,
			Expires:  expirationTime,
			HttpOnly: httpOnly,
			Path:     "/",
		}, nil
	} else {
		return &http.Cookie{
			Name:     name,
			Value:    dataString,
			Expires:  expirationTime,
			HttpOnly: httpOnly,
		}, nil
	}

}
