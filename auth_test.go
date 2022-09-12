package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/thanhpk/randstr"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDBConnect(t *testing.T) {
	dbC := NewDefaultDbCredentials()
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	err = db.Ping()
	if err != nil {
		t.Errorf("[db.Ping] expected nil err, but we got err: %v", err)
	}
	db.Close()
}
func TestDBCreate(t *testing.T) {
	dbC := NewDefaultDbCredentials()
	err := dbC.DbCreateTables()
	if err == nil {
		dbC.DbCreateUsers()
	}

}
func TestDBReqest(t *testing.T) {
	dbC := NewDefaultDbCredentials()
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	var dbCreds []Credentials
	err = db.Select(&dbCreds, `SELECT * FROM users WHERE "password" = $1 and "userName" = $2;`, "admin", "admin")
	if err != nil {
		t.Errorf("[db.Select] expected nil err, but we got err: %v", err)
	}
	db.Close()

}
func TestCreateJWTEmpty(t *testing.T) {
	data, err := CreateJWT("", time.Now(), true, &Claims{}, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	dataExcepted := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{})
	dataStringExcepted, err := dataExcepted.SignedString(jwtKey)
	if err != nil {
		t.Errorf("[jwt.NewWithClaims] expected nil err, but we got err: %v", err)
	}
	if dataStringExcepted != data.Value {
		t.Errorf("expected %v\nbut we got: %v", dataStringExcepted, data.Value)
	}
}
func TestCreateJWThttpOnlyTrue(t *testing.T) {
	httpOnlyBool := true
	data, err := CreateJWT("", time.Now(), httpOnlyBool, &Claims{}, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	if data.HttpOnly != httpOnlyBool {
		t.Errorf("expected data.HttpOnly %v, but we got: %v", httpOnlyBool, data.HttpOnly)
	}
}
func TestCreateJWThttpOnlyFalse(t *testing.T) {
	httpOnlyBool := false
	data, err := CreateJWT("", time.Now(), httpOnlyBool, &Claims{}, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	if data.HttpOnly != httpOnlyBool {
		t.Errorf("expected data.HttpOnly %v, but we got: %v", httpOnlyBool, data.HttpOnly)
	}
}
func TestCreateJWTPathDefault(t *testing.T) {
	exceptedPath := ""
	data, err := CreateJWT("", time.Now(), false, &Claims{}, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	if data.Path != exceptedPath {
		t.Errorf(`expected data.Path "%v", but we got: "%v"`, exceptedPath, data.Path)
	}
}
func TestCreateJWTPathShort(t *testing.T) {
	exceptedPath := "/"
	data, err := CreateJWT("", time.Now(), false, &Claims{}, true)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	if data.Path != exceptedPath {
		t.Errorf(`expected data.Path "%v", but we got: "%v"`, exceptedPath, data.Path)
	}
}
func TestSigninEmpty(t *testing.T) {
	exceptedCode := http.StatusBadRequest
	req := httptest.NewRequest(http.MethodPost, "/auth/signin", nil)
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Signin(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}

func TestSigninWrongLogin(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	reader := strings.NewReader(`{
  "username":"user2",
  "password":"password"
}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Refresh(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestSigninWrongPassword(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	reader := strings.NewReader(`{
  "username":"user2",
  "password":"password"
}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Signin(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}

func TestSigninCorrect(t *testing.T) {
	reader := strings.NewReader(`{
	 "username":"admin",
	 "password":"admin"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Signin(w, req)
	if w.Header()["Set-Cookie"] == nil {
		t.Errorf("expected Set-Cookie in Header, but we got nil")
	}
	if len(w.Header()["Set-Cookie"]) != 3 {
		t.Errorf("excepted Set-Cookie with len 3, but got %v", len(w.Header()["Set-Cookie"]))
	} else {
		if strings.Split(w.Header()["Set-Cookie"][0], "=")[0] != "access_token" {
			t.Errorf("expected access_token, but we got %v", strings.Split(w.Header()["Set-Cookie"][0], "=")[0])
		}
		if strings.Split(w.Header()["Set-Cookie"][1], "=")[0] != "fingerprint" {
			t.Errorf("expected access_token, but we got %v", strings.Split(w.Header()["Set-Cookie"][1], "=")[0])
		}
		if strings.Split(w.Header()["Set-Cookie"][2], "=")[0] != "refresh_token" {
			t.Errorf("expected access_token, but we got %v", strings.Split(w.Header()["Set-Cookie"][2], "=")[0])
		}
		//удаляем созданное подлючение
		db, err := sqlx.Open("postgres", dbC.dbDataSource())
		defer db.Close()
		if err != nil {
			t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
		}
		_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, strings.Split(strings.Split(w.Header()["Set-Cookie"][2], "=")[1], ";")[0])
		if err != nil {
			t.Errorf("[sqlx.Exec] expected nil err, but we got err: %v", err)
		}
	}

}
func TestWelcomeEmpty(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/welcome", nil)
	w := httptest.NewRecorder()
	Welcome(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestWelcomeWrong(t *testing.T) {
	exceptedCode := http.StatusBadRequest
	req := httptest.NewRequest(http.MethodGet, "/welcome", nil)
	expirationTime := time.Now().Add(time.Minute * 11)
	claims := &Claims{
		Username: "username",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := CreateJWT("access_token", expirationTime, false, claims, true)
	accessTokenCookie.Value = accessTokenCookie.Value[0 : len(accessTokenCookie.Value)-8]
	if err != nil {
		t.Errorf("expected accessToken, but we got err: %v", err)
	}
	req.AddCookie(accessTokenCookie)
	w := httptest.NewRecorder()
	Welcome(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestWelcomeWrongOld(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/welcome", nil)
	expirationTime := time.Now().Add(-time.Minute * 11)
	claims := &Claims{
		Username: "username",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := CreateJWT("access_token", expirationTime, false, claims, true)
	if err != nil {
		t.Errorf("expected accessToken, but we got err: %v", err)
	}
	req.AddCookie(accessTokenCookie)
	w := httptest.NewRecorder()
	Welcome(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestWelcomeCorrect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/welcome", nil)
	expirationTime := time.Now().Add(time.Minute * 11)
	claims := &Claims{
		Username: "username",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := CreateJWT("access_token", expirationTime, false, claims, true)
	if err != nil {
		t.Errorf("expected accessToken, but we got err: %v", err)
	}
	req.AddCookie(accessTokenCookie)
	w := httptest.NewRecorder()
	Welcome(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200 code, but we got %v", w.Code)
	}
}
func TestRefreshEmpty(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Refresh(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestRefreshWrongRefreshToken(t *testing.T) {
	exceptedCode := http.StatusBadRequest
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	expirationTime := time.Now().Add(time.Minute * 11)
	claims := &Claims{
		Username: "username",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, false, claims, true)
	refreshTokenCookie.Value = refreshTokenCookie.Value[0 : len(refreshTokenCookie.Value)-8]
	if err != nil {
		t.Errorf("expected accessToken, but we got err: %v", err)
	}
	req.AddCookie(refreshTokenCookie)
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Refresh(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestRefreshTokenNotInDB(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	expirationTime := time.Now().Add(time.Minute * 11)
	claims := &Claims{
		Username: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, false, claims, true)
	if err != nil {
		t.Errorf("expected accessToken, but we got err: %v", err)
	}
	req.AddCookie(refreshTokenCookie)
	w := httptest.NewRecorder()
	dbC := NewDefaultDbCredentials()
	dbC.Refresh(w, req)
	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
func TestRefreshCorrect(t *testing.T) {
	exceptedCode := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	dbC := NewDefaultDbCredentials()
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	var dbCreds []Credentials
	err = db.Select(&dbCreds, `SELECT * FROM users WHERE "password" = $1 and "userName" = $2;`, "password", "user")
	if err != nil {
		t.Errorf("[db.Select] expected nil err, but we got err: %v", err)
	}
	if len(dbCreds) == 0 {
		t.Errorf("[len(dbCreds)] expected >0, but we got: %v", len(dbCreds))
	}

	//создание refreshToken
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		Username: dbCreds[0].Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, true, claims, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	//создание fingerprint
	expirationTimeF := time.Now().Add(time.Hour * 24 * 365)
	fp := &Fp{
		Fingerprint: randstr.Hex(16),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeF.Unix(),
		},
	}
	fingerprintCookie, err := CreateJWT("fingerprint", expirationTimeF, true, fp, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}

	//добавляем в бд данные о токене и устройтсве
	refreshSession := &RefreshSession{
		UserId:      dbCreds[0].UserId,
		ReToken:     refreshTokenCookie.Value,
		UserAgent:   req.UserAgent(),
		Fingerprint: fp.Fingerprint,
		Ip:          strings.Split(req.RemoteAddr, ":")[0],
		ExpiresIn:   expirationTime.Unix(),
	}
	_, err = db.Exec(`insert into refreshSessions ("userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn") values ($1, $2, $3, $4, $5, $6);`,
		refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint, refreshSession.Ip, refreshSession.ExpiresIn)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	req.AddCookie(fingerprintCookie)
	req.AddCookie(refreshTokenCookie)
	w := httptest.NewRecorder()
	dbC.Refresh(w, req)

	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
	_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, refreshSession.ReToken)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}
}

func TestRefreshWrongIp(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	dbC := NewDefaultDbCredentials()
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	var dbCreds []Credentials
	err = db.Select(&dbCreds, `SELECT * FROM users WHERE "password" = $1 and "userName" = $2;`, "password", "user")
	if err != nil {
		t.Errorf("[db.Select] expected nil err, but we got err: %v", err)
	}
	if len(dbCreds) == 0 {
		t.Errorf("[len(dbCreds)] expected >0, but we got: %v", len(dbCreds))
	}

	//создание refreshToken
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		Username: dbCreds[0].Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, true, claims, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}

	//создание fingerprint
	expirationTimeF := time.Now().Add(time.Hour * 24 * 365)
	fp := &Fp{
		Fingerprint: randstr.Hex(16),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeF.Unix(),
		},
	}
	fingerprintCookie, err := CreateJWT("fingerprint", expirationTimeF, true, fp, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}

	//добавляем в бд данные о токене и устройтсве
	refreshSession := &RefreshSession{
		UserId:      dbCreds[0].UserId,
		ReToken:     refreshTokenCookie.Value,
		UserAgent:   req.UserAgent(),
		Fingerprint: fp.Fingerprint,
		Ip:          "0",
		ExpiresIn:   expirationTime.Unix(),
	}
	_, err = db.Exec(`insert into refreshSessions ("userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn") values ($1, $2, $3, $4, $5, $6);`,
		refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint, refreshSession.Ip, refreshSession.ExpiresIn)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	req.AddCookie(fingerprintCookie)
	req.AddCookie(refreshTokenCookie)
	w := httptest.NewRecorder()
	dbC.Refresh(w, req)

	_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, refreshSession.ReToken)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}

func TestRefreshWrongUserAgent(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	dbC := NewDefaultDbCredentials()
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	var dbCreds []Credentials
	err = db.Select(&dbCreds, `SELECT * FROM users WHERE "password" = $1 and "userName" = $2;`, "password", "user")
	if err != nil {
		t.Errorf("[db.Select] expected nil err, but we got err: %v", err)
	}
	if len(dbCreds) == 0 {
		t.Errorf("[len(dbCreds)] expected >0, but we got: %v", len(dbCreds))
	}

	//создание refreshToken
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		Username: dbCreds[0].Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, true, claims, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}

	//создание fingerprint
	expirationTimeF := time.Now().Add(time.Hour * 24 * 365)
	fp := &Fp{
		Fingerprint: randstr.Hex(16),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeF.Unix(),
		},
	}
	fingerprintCookie, err := CreateJWT("fingerprint", expirationTimeF, true, fp, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}

	//добавляем в бд данные о токене и устройтсве
	refreshSession := &RefreshSession{
		UserId:      dbCreds[0].UserId,
		ReToken:     refreshTokenCookie.Value,
		UserAgent:   "Wrong",
		Fingerprint: fp.Fingerprint,
		Ip:          strings.Split(req.RemoteAddr, ":")[0],
		ExpiresIn:   expirationTime.Unix(),
	}
	_, err = db.Exec(`insert into refreshSessions ("userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn") values ($1, $2, $3, $4, $5, $6);`,
		refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint, refreshSession.Ip, refreshSession.ExpiresIn)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	req.AddCookie(fingerprintCookie)
	req.AddCookie(refreshTokenCookie)
	w := httptest.NewRecorder()
	dbC.Refresh(w, req)

	_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, refreshSession.ReToken)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}

func TestRefreshWrongFingerprint(t *testing.T) {
	exceptedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
	dbC := NewDefaultDbCredentials()
	db, err := sqlx.Open("postgres", dbC.dbDataSource())
	defer db.Close()
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	var dbCreds []Credentials
	err = db.Select(&dbCreds, `SELECT * FROM users WHERE "password" = $1 and "userName" = $2;`, "password", "user")
	if err != nil {
		t.Errorf("[db.Select] expected nil err, but we got err: %v", err)
	}
	if len(dbCreds) == 0 {
		t.Errorf("[len(dbCreds)] expected >0, but we got: %v", len(dbCreds))
	}

	//создание refreshToken
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		Username: dbCreds[0].Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := CreateJWT("refresh_token", expirationTime, true, claims, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}

	//создание fingerprint
	expirationTimeF := time.Now().Add(time.Hour * 24 * 365)
	fp := &Fp{
		Fingerprint: randstr.Hex(16),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeF.Unix(),
		},
	}
	fingerprintCookie, err := CreateJWT("fingerprint", expirationTimeF, true, fp, false)
	if err != nil {
		t.Errorf("[CreateJWT] expected nil err, but we got err: %v", err)
	}
	fp.Fingerprint = randstr.Hex(16)
	//добавляем в бд данные о токене и устройтсве
	refreshSession := &RefreshSession{
		UserId:      dbCreds[0].UserId,
		ReToken:     refreshTokenCookie.Value,
		UserAgent:   req.UserAgent(),
		Fingerprint: fp.Fingerprint,
		Ip:          strings.Split(req.RemoteAddr, ":")[0],
		ExpiresIn:   expirationTime.Unix(),
	}
	_, err = db.Exec(`insert into refreshSessions ("userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn") values ($1, $2, $3, $4, $5, $6);`,
		refreshSession.UserId, refreshSession.ReToken, refreshSession.UserAgent, refreshSession.Fingerprint, refreshSession.Ip, refreshSession.ExpiresIn)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	req.AddCookie(fingerprintCookie)
	req.AddCookie(refreshTokenCookie)
	w := httptest.NewRecorder()
	dbC.Refresh(w, req)

	_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, refreshSession.ReToken)
	if err != nil {
		t.Errorf("[db.Exec] expected nil err, but we got err: %v", err)
	}

	if w.Code != exceptedCode {
		t.Errorf("expected %v code, but we got %v", exceptedCode, w.Code)
	}
}
