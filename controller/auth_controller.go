package controller

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/services"
	"github.com/xegcrbq/auth/tokenizer"
	"net/http"
	"time"
)

type AuthController struct {
	service *services.Service
	tknz    *tokenizer.Tokenizer
}

func NewAuthController(service *services.Service, tknz *tokenizer.Tokenizer) *AuthController {
	return &AuthController{
		service: service,
		tknz:    tknz,
	}
}

// Signin обработчик авторизации по логину и паролю
func (a AuthController) Signin(c *fiber.Ctx) error {

	creds := &models.Credentials{
		Username: c.Params("username"),
		Password: c.Params("password"),
	}
	//подключение к бд
	answ := a.service.Execute(models.QueryReadCredentialsByUsername{Username: creds.Username})
	if answ.Err != nil {
		if answ.Err == services.ErrDataNotFound || answ.Err == sql.ErrNoRows {
			c.SendStatus(http.StatusUnauthorized)
			return nil
		}
		return answ.Err
	}
	if answ.Credentials.Password != creds.Password {
		c.SendStatus(http.StatusUnauthorized)
		return nil
	}
	//создание accessToken
	atCookie, err := a.tknz.NewJWTCookie("access_token", creds.Username, time.Now().Add(time.Minute*11))
	if err != nil {
		return err
	}

	//создание fingerprint
	fingerprint := randstr.Hex(16)
	fpCookie, err := a.tknz.NewJWTCookieHTTPOnly("fingerprint", fingerprint, time.Now().Add(time.Hour*24*365))
	if err != nil {
		return err
	}

	//создание refreshToken
	rtCookie, err := a.tknz.NewJWTCookieHTTPOnly("refresh_token", creds.Username, time.Now().Add(time.Hour*24))
	if err != nil {
		return err
	}

	//создание записи в бд
	refreshSession := &models.Session{
		UserId:      creds.UserId,
		ReToken:     rtCookie.Value,
		UserAgent:   c.Get("User-Agent"),
		Fingerprint: fingerprint,
		Ip:          c.IP(),
		ExpiresIn:   rtCookie.Expires.Unix(),
	}
	answ = a.service.Execute(models.CommandCreateSession{Session: refreshSession})
	if answ.Err != nil {
		return err
	}
	c.Cookie(atCookie)
	c.Cookie(fpCookie)
	c.Cookie(rtCookie)
	return nil
}

// Welcome авторизация через cookie
func (a AuthController) Welcome(c *fiber.Ctx) error {
	//получаем cookie
	tokenString := c.Cookies("access_token")
	claims, tkn, err := a.tknz.ParseDataClaims(tokenString)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.SendStatus(http.StatusUnauthorized)
			return nil
		}
		c.SendStatus(http.StatusBadRequest)
		return nil
	}
	if !tkn.Valid {
		c.SendStatus(http.StatusUnauthorized)
		return nil
	}
	c.SendString(fmt.Sprintf("Welcome! %v", claims.Data))
	return nil
}

func (a AuthController) Refresh(c *fiber.Ctx) error {
	//получаем refresh_token cookie
	tokenString := c.Cookies("refresh_token")
	rtClaims, rTkn, err := a.tknz.ParseDataClaims(tokenString)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.SendStatus(http.StatusUnauthorized)
			return err
		}
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	if !rTkn.Valid {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}

	//проверка на наличие токена в бд
	answ := a.service.Execute(models.QueryReadSessionByRefreshToken{RefreshToken: tokenString})
	if answ.Err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	session := answ.Session
	//удаляем из бд сессию с токеном(потом если что вернём, но если не пройдет проверку, то удалится навсегда)
	answ = a.service.Execute(models.CommandDeleteSessionByRefreshToken{tokenString})
	if answ.Err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	//делаем проверку сессии на соотвестсвие данным из бд
	if session.Ip != c.IP() {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}
	if session.UserAgent != c.Get("User-Agent") {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}
	if session.ExpiresIn != rtClaims.ExpiresAt {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}

	//если прошлые проверки пройдены читаем куки с fingerprint
	fpTokenString := c.Cookies("fingerprint")
	fpClaims, fpTkn, err := a.tknz.ParseDataClaims(fpTokenString)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.SendStatus(http.StatusUnauthorized)
			return nil
		}
		c.SendStatus(http.StatusBadRequest)
		return nil
	}
	if !fpTkn.Valid {
		c.SendStatus(http.StatusUnauthorized)
		return nil
	}
	//сравниваем Fingerprint из cookies с Fingerprint из базы данных
	if fpClaims.Data != session.Fingerprint {
		c.SendStatus(http.StatusUnauthorized)
		return nil
	}

	//все проверки пройдены
	//возвращаем сессию в бд
	answ = a.service.Execute(models.CommandCreateSession{Session: session})
	if answ.Err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	//создаем новый access токен
	atCookie, err := a.tknz.NewJWTCookie("access_token", rtClaims.Data, time.Now().Add(time.Minute*11))
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	c.Cookie(atCookie)
	return nil
}
