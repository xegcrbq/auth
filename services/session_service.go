package services

import (
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/repositories"
)

type SessionService struct {
	sessionRepo repositories.SessionRepo
}

func NewSessionService(sessionRepo repositories.SessionRepo) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) IsSessionAvailable(cmd models.CommandDeleteSessionByRefreshToken) (bool, error) {
	session, err := s.sessionRepo.ReadSessionByRefreshToken(&models.QueryReadSessionByRefreshToken{RefreshToken: cmd.RefreshToken})
	if err == nil && session != nil {
		return true, nil
	}
	return false, err
}
func (s *SessionService) GetSession(cmd models.QueryReadSessionByRefreshToken) (*models.Session, error) {
	session, err := s.sessionRepo.ReadSessionByRefreshToken(&cmd)
	return session, err
}
func (s *SessionService) InsertSession(cmd models.CommandCreateSession) error {
	err := s.sessionRepo.SaveSession(&cmd)
	return err
}
func (s *SessionService) DeleteSession(cmd models.CommandDeleteSessionByRefreshToken) error {
	found, err := s.IsSessionAvailable(cmd)
	if !found {
		return ErrDataNotFound
	}
	if err != nil {
		return err
	}
	err = s.sessionRepo.DeleteSessionByRefreshToken(&cmd)
	return err
}
