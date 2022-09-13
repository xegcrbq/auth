package service

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/old"
	"github.com/xegcrbq/auth/old/task"
)

type Service struct {
	rsRepo model.RefreshSessionRepository
	crRepo model.CredentialsRepository
}

func NewService()
func (s *Service) AddRepo(rt old.RepositoryType, db *sqlx.DB) error {
	if db == nil {
		return errors.New("[Service.AddRepo] nil db input")
	}
	if s.repositories == nil {
		s.repositories = make(map[string]model.Repository)
	}
	switch rt {
	case old.REFRESHSESSION:
		s.repositories["RefreshSession"] = old.NewRepo(db)
		return nil
	default:
		return errors.New("[Service.AddRepo] unknown repository type")
	}
}

func (s *Service) RunTask(t task.Task) (model.Model, error) {
	switch t.Model.(type) {
	case model.RefreshSession:
		r, ok := s.repositories["RefreshSession"]
		if !ok {
			return nil, errors.New("[Service.RunTask] repository not found")
		}
		return r.RunTask(t)
	default:
		return nil, errors.New("[Service.RunTask] unknown model type")
	}
}
