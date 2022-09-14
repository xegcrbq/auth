package services

import (
	"github.com/xegcrbq/auth/models"
	"github.com/xegcrbq/auth/repositories"
)

type CredentialsService struct {
	credentialsRepo *repositories.CredentialsRepo
}

func NewCredentialsService(credentialsRepo *repositories.CredentialsRepo) *CredentialsService {
	return &CredentialsService{
		credentialsRepo: credentialsRepo,
	}
}

func (s *CredentialsService) IsUserAvailable(username string) (bool, error) {
	session, err := s.credentialsRepo.ReadCredentialsByUsername(username)
	if err == nil && session != nil {
		return true, nil
	}
	return false, err
}
func (s *CredentialsService) GetCredentials(username string) (*models.Credentials, error) {
	session, err := s.credentialsRepo.ReadCredentialsByUsername(username)
	return session, err
}
func (s *CredentialsService) InsertCredentials(session *models.Credentials) (bool, error) {
	err := s.credentialsRepo.SaveCredentials(session)
	if err != nil {
		return false, err
	}
	return true, err
}
func (s *CredentialsService) DeleteCredentials(username string) (bool, error) {
	found, err := s.IsUserAvailable(username)
	if !found {
		return false, err
	}
	err = s.credentialsRepo.DeleteCredentialsByUsername(username)
	if err != nil {
		return false, err
	}
	return true, err
}
