package repositories

import (
	"fmt"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/models"
	"testing"
	"time"
)

var sr = NewSessionRepo(db.ConnectDB())
var expectedRS = models.Session{
	UserId:      1,
	ReToken:     "TestCRUD" + randstr.Hex(6),
	UserAgent:   "TestCRUD",
	Fingerprint: "TestCRUD",
	Ip:          "TestCRUD",
	ExpiresIn:   time.Now().Add(10 * time.Minute).Unix(),
	CreatedAt:   time.Now(),
}

func TestSessionRepo_SaveSession(t *testing.T) {
	fmt.Println(sr.SaveSession(&expectedRS))
}
func TestSessionRepo_ReadSessionByRefreshToken(t *testing.T) {
	fmt.Println(sr.ReadSessionByRefreshToken(expectedRS.ReToken))
}
func TestSessionRepo_DeleteSessionByRefreshToken(t *testing.T) {
	fmt.Println(sr.DeleteSessionByRefreshToken(expectedRS.ReToken))
}
