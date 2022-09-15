package repositories

import (
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/models"
	"testing"
)

func TestCredentialsRepo(t *testing.T) {
	cr := NewCredentialsRepo(db.ConnectDB())
	expectedCreds := &models.Credentials{
		Username: "TestCRUD" + randstr.Hex(6),
		Password: "TestCRUD" + randstr.Hex(6),
	}
	{
		testID := 0
		t.Logf("\tTest %d:\tSaveCredentials", testID)
		{
			err := cr.SaveCredentials(&models.CommandCreateCredentials{Credentials: expectedCreds})
			assert.Equal(t, nil, err, "expected nil err, but we got: ", err)
		}
		testID++
		t.Logf("\tTest %d:\tReadCredentialsByUsername", testID)
		{
			rCreds, err := cr.ReadCredentialsByUsername(&models.QueryReadCredentialsByUsername{Username: expectedCreds.Username})
			expectedCreds.UserId = rCreds.UserId
			assert.Equal(t, expectedCreds, rCreds, "read incorrect credentials")
			assert.Equal(t, nil, err, "expected nil err, but we got: ", err)
		}
		testID++
		t.Logf("\tTest %d:\tDeleteCredentialsByUsername", testID)
		{
			err := cr.DeleteCredentialsByUsername(&models.CommandDeleteCredentialsByUsername{Username: expectedCreds.Username})
			assert.Equal(t, nil, err, "expected nil err, but we got: ", err)
		}
	}
}
