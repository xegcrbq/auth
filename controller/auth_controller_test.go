package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/repositories"
	"github.com/xegcrbq/auth/services"
	"github.com/xegcrbq/auth/tokenizer"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthController(t *testing.T) {
	sr := repositories.NewSessionRepo(db.ConnectDB())
	cr := repositories.NewCredentialsRepo(db.ConnectDB())
	ss := services.NewSessionService(sr)
	cs := services.NewCredentialsService(cr)
	service := services.NewService(cs, ss)
	tknz := tokenizer.NewTestTokenizer()
	wrongTknz := tokenizer.NewTokenizer([]byte("TestAuthController"))
	a := NewAuthController(service, tknz)
	app := fiber.New()

	app.Get("/auth/:username-:password", a.Signin)
	app.Get("/", a.Welcome)
	app.Get("/auth/refresh", a.Refresh)
	testID := 0
	creds := &models.Credentials{
		Password: "TestAuth",
		Username: "TestAuth" + randstr.Hex(4),
	}
	t.Log("testing SignIn")
	{
		service.Execute(models.CommandCreateCredentials{creds})
		loginTests := []struct {
			username       string
			password       string
			expectedCode   int
			expectedCookie []string
			description    string
		}{
			{
				username:     creds.Username,
				password:     creds.Password,
				expectedCode: http.StatusOK,
				expectedCookie: []string{
					"access_token",
					"fingerprint",
					"refresh_token",
				},
				description: "correct login",
			},
			{
				username:     creds.Username,
				password:     randstr.Hex(10),
				expectedCode: http.StatusUnauthorized,
				description:  "wrong password",
			},
			{
				username:     randstr.Hex(10),
				password:     creds.Username,
				expectedCode: http.StatusUnauthorized,
				description:  "wrong login",
			},
		}

		for _, test := range loginTests {
			t.Logf("\tTest %d:\t%v", testID, test.description)
			testID++
			req := httptest.NewRequest("GET", fmt.Sprintf("/auth/%v-%v", test.username, test.password), nil)
			resp, _ := app.Test(req, 2000)
			assert.Equal(t, test.expectedCode, resp.StatusCode)
			for _, c := range resp.Cookies() {
				assert.Contains(t, test.expectedCookie, c.Name)
			}
		}
	}
	t.Log("testing Welcome")
	{
		correctCookie, _ := tknz.NewJWTCookie("access_token", creds.Username, time.Now().Add(time.Minute))
		wrongCookie, _ := wrongTknz.NewJWTCookie("access_token", creds.Username, time.Now().Add(time.Minute))
		damagedCookie, _ := tknz.NewJWTCookie("access_token", creds.Username, time.Now().Add(time.Minute))
		oldCookie, _ := tknz.NewJWTCookie("access_token", creds.Username, time.Now().Add(-time.Minute))
		damagedCookie.Value = damagedCookie.Value[:len(damagedCookie.Value)-10]
		welcomeTests := []struct {
			cookie       *fiber.Cookie
			expectedCode int
			description  string
		}{
			{
				cookie:       correctCookie,
				expectedCode: http.StatusOK,
				description:  "correct cookie",
			},
			{
				cookie:       wrongCookie,
				expectedCode: http.StatusBadRequest,
				description:  "wrong cookie",
			},
			{
				cookie:       damagedCookie,
				expectedCode: http.StatusBadRequest,
				description:  "damaged cookie",
			},
			{
				cookie:       oldCookie,
				expectedCode: http.StatusBadRequest,
				description:  "old cookie",
			},
		}
		for _, test := range welcomeTests {
			t.Logf("\tTest %d:\t%v", testID, test.description)
			testID++
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(fiberToHttpCookie(test.cookie))
			resp, _ := app.Test(req, 2000)
			assert.Equal(t, test.expectedCode, resp.StatusCode)
		}
	}
}
func fiberToHttpCookie(fc *fiber.Cookie) *http.Cookie {
	return &http.Cookie{
		Name:     fc.Name,
		Value:    fc.Value,
		Path:     fc.Path,
		Domain:   fc.Domain,
		Expires:  fc.Expires,
		HttpOnly: fc.HTTPOnly,
	}
}
