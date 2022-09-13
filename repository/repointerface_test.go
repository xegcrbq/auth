package repository

import (
	"encoding/json"
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
	rsR.Save(&expectedRS)
	data, _ := rsR.ReadByRefreshToken(expectedRS.ReToken)
	fmt.Println(data[0])
	rsR.DeleteByRefreshToken(expectedRS.ReToken)
}
func TestRunInquiryCreate(t *testing.T) {
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
	data, err := json.Marshal(expectedRS)
	fmt.Println(data, err)
	fmt.Println(rsR.RunInquiry(model.RepositoryInquiry{
		Data:                  data,
		RepositoryInquiryCode: model.CREATE,
	}))
}
