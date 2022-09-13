package repository

import (
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/task"
)

type Repository interface {
	create(m model.Model) (model.Model, error)
	read(m model.Model) (model.Model, error)
	delete(m model.Model) (model.Model, error)
	RunTask(t task.Task) (model.Model, error)
}
