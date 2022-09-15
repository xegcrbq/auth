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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthController(t *testing.T) {
	sr := repositories.NewSessionRepo(db.ConnectDB())
	cr := repositories.NewCredentialsRepo(db.ConnectDB())
	ss := services.NewSessionService(sr)
	cs := services.NewCredentialsService(cr)
	service := services.NewService(cs, ss)
	a := NewAuthController(service, []byte("djkhgkjdfgndkjnkdjnvkjkdgkjd"))
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

	}
}

//func TestSigninCorrect(t *testing.T) {
//
//	reader := strings.NewReader(`{
//	 "username":"admin",
//	 "password":"admin"
//	}`)
//
//	req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//	sr := repositories.NewSessionRepo(db.ConnectDB())
//	cr := repositories.NewCredentialsRepo(db.ConnectDB())
//	ss := services.NewSessionService(sr)
//	cs := services.NewCredentialsService(cr)
//	a := NewAuthController(ss, cs, []byte("djkhgkjdfgndkjnkdjnvkjkdgkjd"))
//	a.Signin(w, req)
//	if w.Header()["Set-Cookie"] == nil {
//		t.Errorf("expected Set-Cookie in Header, but we got nil")
//	}
//	if len(w.Header()["Set-Cookie"]) != 3 {
//		t.Errorf("excepted Set-Cookie with len 3, but got %v", len(w.Header()["Set-Cookie"]))
//	} else {
//		if strings.Split(w.Header()["Set-Cookie"][0], "=")[0] != "access_token" {
//			t.Errorf("expected access_token, but we got %v", strings.Split(w.Header()["Set-Cookie"][0], "=")[0])
//		}
//		if strings.Split(w.Header()["Set-Cookie"][1], "=")[0] != "fingerprint" {
//			t.Errorf("expected access_token, but we got %v", strings.Split(w.Header()["Set-Cookie"][1], "=")[0])
//		}
//		if strings.Split(w.Header()["Set-Cookie"][2], "=")[0] != "refresh_token" {
//			t.Errorf("expected access_token, but we got %v", strings.Split(w.Header()["Set-Cookie"][2], "=")[0])
//		}
//		//удаляем созданное подлючение
//		db, err := sqlx.Open("postgres", dbC.dbDataSource())
//		defer db.Close()
//		if err != nil {
//			t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
//		}
//		_, err = db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, strings.Split(strings.Split(w.Header()["Set-Cookie"][2], "=")[1], ";")[0])
//		if err != nil {
//			t.Errorf("[sqlx.Exec] expected nil err, but we got err: %v", err)
//		}
//	}
//
//}
