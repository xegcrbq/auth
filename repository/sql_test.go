package repository

import (
	"context"
	"fmt"
	"github.com/thanhpk/randstr"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/test"
	"testing"
	"time"
)

func TestNewSqlRefreshSessionRepository(t *testing.T) {
	rsR := NewSqlRefreshSessionRepository(test.Db())
	expectedRS := model.RefreshSession{
		UserId:      1,
		ReToken:     "TestCreate" + randstr.Hex(6),
		UserAgent:   "TestCreate",
		Fingerprint: "TestCreate",
		Ip:          "TestCreate",
		ExpiresIn:   time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:   time.Now(),
	}
	rsR.Save(context.Background(), &expectedRS)
	data, _ := rsR.ReadByRefreshToken(context.Background(), expectedRS.ReToken)
	fmt.Println(data[0])
	rsR.DeleteByRefreshToken(context.Background(), expectedRS.ReToken)
}
