package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/services"
	"net/http"
	"time"
)

type AuthController struct {
	sessionService     *services.SessionService
	credentialsService *services.CredentialsService
	jwtKey             []byte
}

func NewAuthController(sessionService *services.SessionService, credentialsService *services.CredentialsService, jwtKey []byte) *AuthController {
	return &AuthController{
		sessionService:     sessionService,
		credentialsService: credentialsService,
		jwtKey:             jwtKey,
	}
}

// Signin обработчик авторизации по логину и паролю
func (a AuthController) Signin(c *fiber.Ctx) error {

	creds := &models.Credentials{
		Username: c.Params("username"),
		Password: c.Params("password"),
	}
	//подключение к бд
	dbCreds, err := a.credentialsService.GetCredentials(creds.Username)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	//создание accessToken
	expirationTime := time.Now().Add(time.Minute * 11)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := models.CreateJWT("access_token", expirationTime, false, claims, true, a.jwtKey)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	//создание fingerprint
	expirationTime = time.Now().Add(time.Hour * 24 * 365)
	fp := &models.Fp{
		Fingerprint: randstr.Hex(16),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	fingerprintCookie, err := models.CreateJWT("fingerprint", expirationTime, true, fp, false, a.jwtKey)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	//создание refreshToken
	expirationTime = time.Now().Add(time.Hour * 24)

	claims = &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshTokenCookie, err := models.CreateJWT("refresh_token", expirationTime, true, claims, false, a.jwtKey)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	//создание записи в бд
	refreshSession := &models.Session{
		UserId:      dbCreds.UserId,
		ReToken:     refreshTokenCookie.Value,
		UserAgent:   c.Get("User-Agent"),
		Fingerprint: fp.Fingerprint,
		Ip:          c.IP(),
		ExpiresIn:   expirationTime.Unix(),
	}
	_, err = a.sessionService.InsertSession(refreshSession)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	c.Cookie(accessTokenCookie)
	c.Cookie(fingerprintCookie)
	c.Cookie(refreshTokenCookie)
	return nil
}

// Welcome авторизация через cookie
func (a AuthController) Welcome(c *fiber.Ctx) error {
	fmt.Println(c.Cookies("access_token"))
	//получаем cookie
	tokenString := c.Cookies("access_token")
	//создаём структуру для парсинга в неё данных
	claims := &models.Claims{}
	//парсим данные из токена
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.SendStatus(http.StatusUnauthorized)
			return err
		}
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	if !tkn.Valid {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}
	c.SendString(fmt.Sprintf("Welcome! %v", claims.Username))
	return nil
}

func (a AuthController) Refresh(c *fiber.Ctx) error {
	//получаем refresh_token cookie
	tokenString := c.Cookies("refresh_token")
	//создаём структуру для парсинга в неё данных
	claims := &models.Claims{}
	//парсим данные из токена
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.SendStatus(http.StatusUnauthorized)
			return err
		}
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	if !tkn.Valid {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}

	//проверка на наличие токена в бд
	dbRefreshSession, err := a.sessionService.GetSession(tokenString)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	//удаляем из бд сессию с токеном(потом если что вернём, но если не пройдет проверку, то удалится навсегда)
	_, err = a.sessionService.DeleteSession(tokenString)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	//делаем проверку сессии на соотвестсвие данным из бд
	if dbRefreshSession.Ip != c.IP() {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}
	if dbRefreshSession.UserAgent != c.Get("User-Agent") {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}
	if dbRefreshSession.ExpiresIn != claims.ExpiresAt {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}

	//если прошлые проверки пройдены читаем куки с fingerprint
	fTokenString := c.Cookies("fingerprint")
	fp := &models.Fp{}
	ftkn, err := jwt.ParseWithClaims(fTokenString, fp, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.SendStatus(http.StatusUnauthorized)
			return err
		}
		c.SendStatus(http.StatusBadRequest)
		return err
	}
	if !ftkn.Valid {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}
	//сравниваем Fingerprint из cookies с Fingerprint из базы данных
	if fp.Fingerprint != dbRefreshSession.Fingerprint {
		c.SendStatus(http.StatusUnauthorized)
		return err
	}

	//все проверки пройдены

	//возвращаем сессию в бд
	_, err = a.sessionService.InsertSession(dbRefreshSession)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	//создаем новый access токен
	expirationTime := time.Now().Add(time.Minute * 11)
	accessClaims := &models.Claims{
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessTokenCookie, err := models.CreateJWT("access_token", expirationTime, false, accessClaims, true, a.jwtKey)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}
	c.Cookie(accessTokenCookie)
	return nil
}
