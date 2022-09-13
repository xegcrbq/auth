package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/xegcrbq/auth/model"
	"github.com/xegcrbq/auth/sqlstore"
)

type ctxTransactionKey struct{}

type SqlRefreshSessionRepository struct {
	db *sqlx.DB
}
type sqlRepository interface {
	getDB() *sqlx.DB
}

var ErrInvalidTxType = errors.New("invalid tx type, tx type should be *sqlx.Tx")

func getSqlxDatabase(ctx context.Context, r sqlRepository) (sqlstore.SqlxDatabase, error) {
	txv := ctx.Value(ctxTransactionKey{})
	if txv == nil {
		return r.getDB(), nil
	}
	if tx, ok := txv.(*sqlx.Tx); ok {
		return tx, nil
	}
	return nil, ErrInvalidTxType
}

func NewSqlRefreshSessionRepository(db *sqlx.DB) *SqlRefreshSessionRepository {
	return &SqlRefreshSessionRepository{db: db}
}

func (r *SqlRefreshSessionRepository) getDB() *sqlx.DB {
	return r.db
}

func (r *SqlRefreshSessionRepository) Save(ctx context.Context, rs *model.RefreshSession) error {
	db, err := getSqlxDatabase(ctx, r)
	if err != nil {
		return err
	}
	return sqlstore.Save(ctx, db, rs)
}

func (r *SqlRefreshSessionRepository) ReadByRefreshToken(ctx context.Context, refreshToken string) ([]*model.RefreshSession, error) {
	db, err := getSqlxDatabase(ctx, r)
	if err != nil {
		return nil, err
	}
	return sqlstore.ReadByRefreshToken(ctx, db, refreshToken)
}

func (r *SqlRefreshSessionRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	db, err := getSqlxDatabase(ctx, r)
	if err != nil {
		return err
	}
	return sqlstore.DeleteByRefreshToken(ctx, db, refreshToken)

}
