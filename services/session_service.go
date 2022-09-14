package services

import (
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/repositories"
)

type SessionService struct {
	sessionRepo *repositories.SessionRepo
}

func NewSessionService(sessionRepo *repositories.SessionRepo) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) IsSessionAvailable(refreshToken string) (bool, error) {
	session, err := s.sessionRepo.ReadSessionByRefreshToken(refreshToken)
	if err == nil && session != nil {
		return true, nil
	}
	return false, err
}
func (s *SessionService) GetSession(refreshToken string) (*models.Session, error) {
	session, err := s.sessionRepo.ReadSessionByRefreshToken(refreshToken)
	return session, err
}
func (s *SessionService) InsertSession(session *models.Session) (bool, error) {
	err := s.sessionRepo.SaveSession(session)
	if err != nil {
		return false, err
	}
	return true, err
}
func (s *SessionService) DeleteSession(refreshToken string) (bool, error) {
	found, err := s.IsSessionAvailable(refreshToken)
	if !found {
		return false, err
	}
	err = s.sessionRepo.DeleteSessionByRefreshToken(refreshToken)
	if err != nil {
		return false, err
	}
	return true, err
}
