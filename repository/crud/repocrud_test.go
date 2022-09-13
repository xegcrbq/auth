package crud

import (
	"context"
	"fmt"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/test"
	"testing"
	"time"
)

func TestDropTables(t *testing.T) {
	DropTables(test.Db())
	CreateTables(test.Db())
}

func TestSaveRefreshSession(t *testing.T) {
	expectedRS := model.RefreshSession{
		UserId:      1,
		ReToken:     "TestCreate" + randstr.Hex(6),
		UserAgent:   "TestCreate",
		Fingerprint: "TestCreate",
		Ip:          "TestCreate",
		ExpiresIn:   time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:   time.Now(),
	}
	Save(context.Background(), test.Db(), &expectedRS)
	data, _ := ReadByRefreshToken(context.Background(), test.Db(), expectedRS.ReToken)
	fmt.Println(data[0])
	DeleteByRefreshToken(context.Background(), test.Db(), expectedRS.ReToken)
}
