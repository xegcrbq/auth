package repositories

import (
	"fmt"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/db"
	"github.com/xegcrbq/auth/models"
	"testing"
)

var cr = NewCredentialsRepo(db.ConnectDB())
var expectedC = models.Credentials{
	Username: "TestCRUD" + randstr.Hex(6),
	Password: "TestCRUD" + randstr.Hex(6),
}

func TestCredentialsRepo_SaveCredentials(t *testing.T) {
	fmt.Println(cr.SaveCredentials(&expectedC))
}
func TestCredentialsRepo_ReadCredentialsByUsername(t *testing.T) {
	fmt.Println(cr.ReadCredentialsByUsername(expectedC.Username))
}
func TestCredentialsRepo_DeleteCredentialsByUsername(t *testing.T) {
	fmt.Println(cr.DeleteCredentialsByUsername(expectedC.Username))
}
