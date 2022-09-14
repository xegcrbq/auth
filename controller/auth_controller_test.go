package controller

//
//import (
//	"github.com/jmoiron/sqlx"
//	"github.com/xegcrbq/auth/db"
//	"github.com/xegcrbq/auth/repositories"
//	"github.com/xegcrbq/auth/services"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func TestSigninCorrect(t *testing.T) {
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
