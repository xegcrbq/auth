package repositories

import (
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/models"
	"testing"
	"time"
)

func TestSessionRepo(t *testing.T) {
	sr := NewSessionRepoSQL(db.ConnectDB())
	expectedSession := &models.Session{
		UserId:      1,
		ReToken:     "TestCRUD" + randstr.Hex(6),
		UserAgent:   "TestCRUD",
		Fingerprint: "TestCRUD",
		Ip:          "TestCRUD",
		ExpiresIn:   time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:   time.Now(),
	}
	{
		testID := 0
		t.Logf("\tTest %d:\tSaveSession", testID)
		{
			err := sr.SaveSession(&models.CommandCreateSession{Session: expectedSession})
			assert.Equal(t, nil, err, "expected nil err, but we got: ", err)
		}
		testID++
		t.Logf("\tTest %d:\tReadSessionByRefreshToken", testID)
		{
			rSession, err := sr.ReadSessionByRefreshToken(&models.QueryReadSessionByRefreshToken{RefreshToken: expectedSession.ReToken})
			expectedSession.Id = rSession.Id
			expectedSession.CreatedAt = expectedSession.CreatedAt.Round(time.Second).Local()
			rSession.CreatedAt = rSession.CreatedAt.Round(time.Second).Local()
			assert.Equal(t, expectedSession, rSession, "read incorrect session")
			assert.Equal(t, nil, err, "expected nil err, but we got: ", err)
		}
		testID++
		t.Logf("\tTest %d:\tDeleteSessionByRefreshToken", testID)
		{
			err := sr.DeleteSessionByRefreshToken(&models.CommandDeleteSessionByRefreshToken{RefreshToken: expectedSession.ReToken})
			assert.Equal(t, nil, err, "expected nil err, but we got: ", err)
		}
	}
}
