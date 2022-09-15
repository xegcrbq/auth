package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/repositories"
	"testing"
	"time"
)

func TestCredentialsRepo(t *testing.T) {
	sr := repositories.NewSessionRepo(db.ConnectDB())
	cr := repositories.NewCredentialsRepo(db.ConnectDB())
	ss := NewSessionService(sr)
	cs := NewCredentialsService(cr)
	service := NewService(cs, ss)
	expectedCreds := &models.Credentials{
		Username: "TestService" + randstr.Hex(6),
		Password: "TestService" + randstr.Hex(6),
	}
	expectedSession := &models.Session{
		UserId:      1,
		ReToken:     "TestService" + randstr.Hex(6),
		UserAgent:   "TestService",
		Fingerprint: "TestService",
		Ip:          "TestService",
		ExpiresIn:   time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:   time.Now(),
	}
	{
		testID := 0
		t.Logf("\tTest %d:\tCommandCreateCredentials", testID)
		{
			answer := service.Execute(models.CommandCreateCredentials{Credentials: expectedCreds})
			assert.Equal(t, models.Answer{}, *answer)
		}
		testID++
		t.Logf("\tTest %d:\tQueryReadCredentialsByUsername", testID)
		{
			answer := service.Execute(models.QueryReadCredentialsByUsername{expectedCreds.Username})
			expectedCreds.UserId = answer.Credentials.UserId
			assert.Equal(t, models.Answer{Credentials: expectedCreds}, *answer)
		}
		testID++
		t.Logf("\tTest %d:\tCommandDeleteCredentialsByUsername", testID)
		{
			answer := service.Execute(models.CommandDeleteCredentialsByUsername{expectedCreds.Username})
			assert.Equal(t, models.Answer{}, *answer)
		}
		testID++
		t.Logf("\tTest %d:\tCommandCreateSession", testID)
		{
			answer := service.Execute(models.CommandCreateSession{Session: expectedSession})
			assert.Equal(t, models.Answer{}, *answer)
		}
		testID++
		t.Logf("\tTest %d:\tQueryReadSessionByRefreshToken", testID)
		{
			answer := service.Execute(models.QueryReadSessionByRefreshToken{expectedSession.ReToken})
			expectedSession.Id = answer.Session.Id
			expectedSession.CreatedAt = expectedSession.CreatedAt.Round(time.Second).Local()
			answer.Session.CreatedAt = answer.Session.CreatedAt.Round(time.Second).Local()
			assert.Equal(t, models.Answer{Session: expectedSession}, *answer)
		}
		testID++
		t.Logf("\tTest %d:\tCommandDeleteSessionByRefreshToken", testID)
		{
			answer := service.Execute(models.CommandDeleteSessionByRefreshToken{expectedSession.ReToken})
			assert.Equal(t, models.Answer{}, *answer)
		}
	}
}
