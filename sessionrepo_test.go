package auth

import (
	"github.com/jmoiron/sqlx"
	"math/rand"
	"testing"
	"time"
)

var reToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJleHAiOjE2NjMwNzI3NzB9.YJIJIZ7Sk5NugdXnxnqrLbCLB8BHFOBwOCavhLnBb9g"

func TestAddData(t *testing.T) {
	db, err := sqlx.Open("postgres", NewDefaultDbCredentials().dbDataSource())
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	defer db.Close()
	sr := NewSessionRepo(db)
	_, err = sr.AddData(RefreshSession{
		Id:          int64(rand.Int() % 100000),
		UserId:      "1",
		ReToken:     reToken,
		UserAgent:   "test",
		Fingerprint: "fingerprint",
		Ip:          "192.168.1.1",
		ExpiresIn:   10,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Errorf("[AddData] expected nil err, but we got err: %v", err)
	}
}
func TestGetExistingData(t *testing.T) {
	db, err := sqlx.Open("postgres", NewDefaultDbCredentials().dbDataSource())
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	defer db.Close()
	sr := NewSessionRepo(db)
	session, err := sr.GetExistingData(reToken)
	if err != nil {
		t.Errorf("expected nil err, but we got err: %v", err)
	} else if session.ReToken != reToken {
		t.Errorf("expected %v token, but we got %v token", reToken, session.ReToken)
	}
}
func TestRemoveData(t *testing.T) {
	db, err := sqlx.Open("postgres", NewDefaultDbCredentials().dbDataSource())
	if err != nil {
		t.Errorf("[sqlx.Open] expected nil err, but we got err: %v", err)
	}
	defer db.Close()
	sr := NewSessionRepo(db)
	session, err := sr.RemoveData(reToken)
	if err != nil {
		t.Errorf("expected nil err, but we got err: %v", err)
	} else if session.ReToken != reToken {
		t.Errorf("expected %v token, but we got %v token", reToken, session.ReToken)
	}
}
