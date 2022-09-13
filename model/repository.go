package model

import "context"

type RefreshSessionRepository interface {
	SaveRefreshSession(ctx context.Context, rs *RefreshSession) error
	ReadByRefreshToken(ctx context.Context, refreshToken string) ([]*RefreshSession, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
}

type CredentialsRepository interface {
	SaveCredentials(ctx context.Context, rs *Credentials) error
	ReadByUsername(ctx context.Context, username string) ([]*Credentials, error)
	DeleteByUsername(ctx context.Context, username string) error
}
