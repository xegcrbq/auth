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

func (s *CredentialsService) IsUserAvailable(cmd models.CommandDeleteCredentialsByUsername) (bool, error) {
	session, err := s.credentialsRepo.ReadCredentialsByUsername(&models.QueryReadCredentialsByUsername{Username: cmd.Username})
	if err == nil && session != nil {
		return true, nil
	}
	return false, err
}
func (s *CredentialsService) GetCredentials(cmd models.QueryReadCredentialsByUsername) (*models.Credentials, error) {
	session, err := s.credentialsRepo.ReadCredentialsByUsername(&cmd)
	return session, err
}
func (s *CredentialsService) InsertCredentials(cmd models.CommandCreateCredentials) error {
	err := s.credentialsRepo.SaveCredentials(&cmd)
	return err
}
func (s *CredentialsService) DeleteCredentials(cmd models.CommandDeleteCredentialsByUsername) error {
	found, err := s.IsUserAvailable(cmd)
	if !found {
		return ErrDataNotFound
	}
	if err != nil {
		return err
	}
	err = s.credentialsRepo.DeleteCredentialsByUsername(&cmd)
	return err
}
