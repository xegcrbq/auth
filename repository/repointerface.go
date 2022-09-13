package repository

import (
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/repository/crud"
)

type SqlRefreshSessionRepository struct {
	db *sqlx.DB
}
type sqlRepository interface {
	getDB() *sqlx.DB
}

var errNotFound = errors.New("inquiry not found by refresh session repository")

func getSqlxDatabase(r sqlRepository) (crud.SqlxDatabase, error) {
	return r.getDB(), nil
}

func NewSqlRefreshSessionRepository(db *sqlx.DB) *SqlRefreshSessionRepository {
	return &SqlRefreshSessionRepository{db: db}
}

func (r *SqlRefreshSessionRepository) getDB() *sqlx.DB {
	return r.db
}

func (r *SqlRefreshSessionRepository) Save(rs *model.RefreshSession) error {
	db, err := getSqlxDatabase(r)
	if err != nil {
		return err
	}
	return crud.Save(db, rs)
}

func (r *SqlRefreshSessionRepository) ReadByRefreshToken(refreshToken string) ([]*model.RefreshSession, error) {
	db, err := getSqlxDatabase(r)
	if err != nil {
		return nil, err
	}
	return crud.ReadByRefreshToken(db, refreshToken)
}

func (r *SqlRefreshSessionRepository) DeleteByRefreshToken(refreshToken string) error {
	db, err := getSqlxDatabase(r)
	if err != nil {
		return err
	}
	return crud.DeleteByRefreshToken(db, refreshToken)
}

func (r *SqlRefreshSessionRepository) RunInquiry(ri model.RepositoryInquiry) model.Answer {
	switch ri.RepositoryInquiryCode {
	case model.CREATE:
		var rs model.RefreshSession
		err := json.Unmarshal(ri.Data, &rs)
		if err != nil {
			return model.Answer{
				Data: []byte(err.Error()),
				Code: model.ERROR,
			}
		}
		err = r.Save(&rs)
		if err != nil {
			return model.Answer{
				Data: []byte(err.Error()),
				Code: model.ERROR,
			}
		}
		return model.Answer{
			Data: nil,
			Code: model.SUCCSESS,
		}
	case model.READ:
		result, err := r.ReadByRefreshToken(string(ri.Data))
		if err != nil {
			return model.Answer{
				Data: []byte(err.Error()),
				Code: model.ERROR,
			}
		}
		byteResult, err := json.Marshal(result)
		if err != nil {
			return model.Answer{
				Data: []byte(err.Error()),
				Code: model.ERROR,
			}
		}
		return model.Answer{
			Data: byteResult,
			Code: model.DATAREFRESHSESSION,
		}
	default:
		return model.Answer{
			Data: []byte(errNotFound.Error()),
			Code: model.ERROR,
		}
	}

}
