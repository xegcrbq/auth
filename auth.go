package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/old/service"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	service service.Service
}

// Signin обработчик авторизации по логину и паролю
func (a Auth) Signin(w http.ResponseWriter, r *http.Request) {

	var creds model.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	//проверка корректности парсинга
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//подключение к бд
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//проверка на наличие логина и пароля в бд
	var dbCreds []model.Credentials
	err = db.Select(&dbCreds, `SELECT * FROM users WHERE "password" = $1 and "userName" = $2;`, creds.Password, creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(dbCreds) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//создание accessToken
	expirationTime := time.Now().Add(time.Minute * 11)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := CreateJWT("access_token", expirationTime, false, claims, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//создание fingerprint
	expirationTime = time.Now().Add(time.Hour * 24 * 365)
	fp := &Fp{
		Fingerprint: randstr.Hex(16),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	fingerprintCookie, err := CreateJWT("fingerprint", expirationTime, true, fp, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//создание refreshToken
	expirationTime = time.Now().Add(time.Hour * 24)

	claims = &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, true, claims, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//создание записи в бд
	refreshSession := &RefreshSession{
		UserId:      dbCreds[0].UserId,
		ReToken:     refreshTokenCookie.Value,
		UserAgent:   r.UserAgent(),
		Fingerprint: fp.Fingerprint,
		Ip:          strings.Split(r.RemoteAddr, ":")[0],
		ExpiresIn:   expirationTime.Unix(),
	}
	_, err = db.Exec(`insert into refreshSessions ("userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn") values ($1, $2, $3, $4, $5, $6);`,
		refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint, refreshSession.Ip, refreshSession.ExpiresIn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, fingerprintCookie)
	http.SetCookie(w, refreshTokenCookie)

}

// Welcome авторизация через cookie
func (a Auth) Welcome(w http.ResponseWriter, r *http.Request) {

	//получаем cookie
	c, err := r.Cookie("access_token")
	if err != nil {
		//cookie не содержат данные о токене
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//получаем значение токена
	tokenString := c.Value
	//создаём структуру для парсинга в неё данных
	claims := &Claims{}
	//парсим данные из токена
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Now().Unix() > claims.ExpiresAt {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome! %v", claims.Username)))
}

func (a Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	//получаем refresh_token cookie
	c, err := r.Cookie("refresh_token")
	if err != nil {
		//cookie не содержат данные о токене
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//получаем значение токена
	tokenString := c.Value
	//создаём структуру для парсинга в неё данных
	claims := &Claims{}
	//парсим данные из токена
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Now().Unix() > claims.ExpiresAt {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//подключение к бд
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//проверка на наличие токена в бд
	var dbRefreshSessions []RefreshSession
	err = db.Select(&dbRefreshSessions, `SELECT * FROM refreshsessions WHERE "refreshToken" = $1;`, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(dbRefreshSessions) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//удаляем из бд сессию с токеном(потом если что вернём, но если не пройдет проверку, то удалится навсегда)
	_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`,
		dbRefreshSessions[0].ReToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//делаем проверку сессии на соотвестсвие данным из бд
	if dbRefreshSessions[0].Ip != strings.Split(r.RemoteAddr, ":")[0] {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if dbRefreshSessions[0].UserAgent != r.UserAgent() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if dbRefreshSessions[0].ExpiresIn != claims.ExpiresAt {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//если прошлые проверки пройдены читаем куки с fingerprint
	fc, err := r.Cookie("fingerprint")
	if err != nil {
		//cookie не содержат данные о токене
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fTokenString := fc.Value
	fp := &Fp{}
	ftkn, err := jwt.ParseWithClaims(fTokenString, fp, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ftkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//сравниваем Fingerprint из cookies с Fingerprint из базы данных
	if fp.Fingerprint != dbRefreshSessions[0].Fingerprint {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//все проверки пройдены

	//повзвращаем сессию в бд
	_, err = db.Exec(`insert into refreshSessions ("id", "userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn", "createdAt") values ($1, $2, $3, $4, $5, $6, $7, $8);`,
		dbRefreshSessions[0].Id, dbRefreshSessions[0].UserId, dbRefreshSessions[0].ReToken, dbRefreshSessions[0].UserAgent, dbRefreshSessions[0].Fingerprint,
		dbRefreshSessions[0].Ip, dbRefreshSessions[0].ExpiresIn, dbRefreshSessions[0].CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//создаем новый access токен
	expirationTime := time.Now().Add(time.Minute * 11)
	accessClaims := &Claims{
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := CreateJWT("access_token", expirationTime, false, accessClaims, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, accessTokenCookie)
}
