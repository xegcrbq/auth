package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/task"
)

type RefreshSessionRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *RefreshSessionRepo {
	return &RefreshSessionRepo{
		db: db,
	}
}

func (r *RefreshSessionRepo) read(m model.Model) (model.Model, error) {
	rs := m.(model.RefreshSession)
	var outputRS model.RefreshSession
	r.db.Get(&outputRS, `SELECT * FROM refreshsessions WHERE "refreshToken" = $1;`, rs.ReToken)
	if outputRS.Id != 0 {
		return outputRS, nil
	}
	return nil, errors.New("[SessionRepo.read] model not found")
}

func (r *RefreshSessionRepo) delete(m model.Model) (model.Model, error) {
	rs := m.(model.RefreshSession)
	outputRS, err := r.read(rs)
	if err != nil {
		if err.Error() == "[SessionRepo.read] model not found" {
			return nil, errors.New("[SessionRepo.delete] model not found")
		}
		return nil, err
	}
	res, err := r.db.Exec(`DELETE FROM refreshsessions WHERE "refreshToken" = $1;`, rs.ReToken)
	if err != nil {
		return nil, err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 0 && err == nil {
		return nil, err
	}
	return outputRS, nil
}
func (r *RefreshSessionRepo) create(m model.Model) (model.Model, error) {
	rs := m.(model.RefreshSession)
	_, err := r.db.Exec(`INSERT INTO refreshSessions ("userId", "refreshToken", "ua", "fingerprint", "ip", "expiresIn", "createdAt") VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		rs.UserId, rs.ReToken, rs.UserAgent, rs.Fingerprint, rs.Ip, rs.ExpiresIn, rs.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
func (r *RefreshSessionRepo) RunTask(t task.Task) (model.Model, error) {
	switch t.TaskType {
	case task.CREATE:
		return r.create(t.Model)
	case task.READ:
		return r.read(t.Model)
	case task.DELETE:
		return r.delete(t.Model)
	default:
		return nil, errors.New("[SessionRepo.RunTask] task not found")
	}
}
