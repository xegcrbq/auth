package old

import (
	"fmt"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/test"
	"testing"
	"time"
)

var reToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJleHAiOjE2NjMwNzI3NzB9.YJIJIZ7Sk5NugdXnxnqrLbCLB8BHFOBwOCavhLnBb9g"
var db = test.Db()
var refreshsession = NewRepo(db)

func TestConnectDB(t *testing.T) {
	if db == nil {
		t.Errorf("test db connection failed")
	}
	err := db.Ping()
	if err != nil {
		t.Errorf("expected nil err, but we got err: %v", err)
	}
}
func TestNewRepo(t *testing.T) {
	if refreshsession == nil {
		t.Errorf("expected not nil value, but we got nil")
	}
}
func TestCreateCorrect(t *testing.T) {
	expectedRS := model.RefreshSession{
		UserId:      1,
		ReToken:     "TestCreate" + randstr.Hex(6),
		UserAgent:   "TestCreate",
		Fingerprint: "TestCreate",
		Ip:          "TestCreate",
		ExpiresIn:   time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:   time.Now(),
	}
	m, err := refreshsession.create(expectedRS)
	if err != nil {
		t.Errorf("expected nil err, but we got err: %v", err)
	}
	var outputRS model.RefreshSession
	db.Get(&outputRS, `SELECT * FROM refreshsessions WHERE "refreshToken" = $1;`, expectedRS.ReToken)
	if !outputRS.IsValid() {
		t.Errorf("not valid created object")
	}
	if !(expectedRS.Equal(m.(model.RefreshSession)) && expectedRS.Equal(outputRS)) {
		t.Errorf("created data not equal expected")
	}
	res, err := db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, expectedRS.ReToken)
	if err != nil {
		t.Errorf("expected nil err, but we got err: %v", err)
	}
	rAffected, _ := res.RowsAffected()
	if rAffected != 1 {
		t.Errorf("expected 1 RowsAffected, but we got: %v", rAffected)
	}
}
func TestBase(t *testing.T) {
	rs := model.RefreshSession{UserAgent: "TestCreate"}
	var outputRS []model.RefreshSession
	refreshsession.db.Select(&outputRS, `SELECT * FROM refreshsessions WHERE "ua" = $1;`, rs.UserAgent)
	fmt.Println(outputRS)
}

//
//func TestGetExistingData(t *testing.T) {
//	db, err := sqlx.Open("postgres", auth.NewDefaultDbCredentials().dbDataSource())
//	if err != nil {
//		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
//	}
//	defer db.Close()
//	sr := NewSessionRepo(db)
//	session, err := sr.GetExistingData(reToken)
//	if err != nil {
//		t.Errorf("expected nil err, but we got err: %v", err)
//	} else if session.ReToken != reToken {
//		t.Errorf("expected %v token, but we got %v token", reToken, session.ReToken)
//	}
//}
//func TestRemoveData(t *testing.T) {
//	db, err := sqlx.Open("postgres", auth.NewDefaultDbCredentials().dbDataSource())
//	if err != nil {
//		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
//	}
//	defer db.Close()
//	sr := NewSessionRepo(db)
//	session, err := sr.RemoveData(reToken)
//	if err != nil {
//		t.Errorf("expected nil err, but we got err: %v", err)
//	} else if session.ReToken != reToken {
//		t.Errorf("expected %v token, but we got %v token", reToken, session.ReToken)
//	}
//}
