package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// ключ для создания подписи
var jwtKey = []byte("irjgdkngfdkjdkjlbvnjkd")

//var dbDataSource = fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
//	"postgres", 1, "postgres")

type Data interface {
	Valid() error
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
