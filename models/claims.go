package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Data interface {
	Valid() error
}

func (c Claims) Valid() error {
	if c.ExpiresAt < time.Now().Unix() {
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

func (f Fp) Valid() error {
	if f.ExpiresAt < time.Now().Unix() {
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

func CreateJWT(name string, expirationTime time.Time, httpOnly bool, rawData Data, shortPath bool, jwtKey []byte) (*fiber.Cookie, error) {

	data := jwt.NewWithClaims(jwt.SigningMethodHS256, rawData)
	dataString, err := data.SignedString(jwtKey)
	if err != nil {
		return &fiber.Cookie{}, errors.New("data.SignedString(jwtKey) error")
	}
	//выставляем параметр HttpOnly, чтобы получать доступ к этому токену только на странице авторизации
	if shortPath {
		return &fiber.Cookie{
			Name:     name,
			Value:    dataString,
			Expires:  expirationTime,
			HTTPOnly: httpOnly,
			Path:     "/",
		}, nil
	} else {
		return &fiber.Cookie{
			Name:     name,
			Value:    dataString,
			Expires:  expirationTime,
			HTTPOnly: httpOnly,
			Path:     "/auth/",
		}, nil
	}

}
