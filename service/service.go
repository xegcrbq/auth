package service

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/repository"
	"github.com/xegcrbq/auth/task"
)

type Service struct {
	repositories map[string]model.Repository
}

func (s *Service) AddRepo(rt repository.RepositoryType, db *sqlx.DB) error {
	if db == nil {
		return errors.New("[Service.AddRepo] nil db input")
	}
	if s.repositories == nil {
		s.repositories = make(map[string]model.Repository)
	}
	switch rt {
	case repository.REFRESHSESSION:
		s.repositories["RefreshSession"] = repository.NewRepo(db)
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
