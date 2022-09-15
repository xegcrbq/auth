package services

import (
	"errors"
	"github.com/xegcrbq/auth/models"
)

var ErrDataNotFound = errors.New("data not found")
var ErrComandNotFound = errors.New("command not found")

type Service struct {
	credsService   *CredentialsService
	sessionService *SessionService
}

func NewService(credsService *CredentialsService, sessionService *SessionService) *Service {
	return &Service{
		credsService:   credsService,
		sessionService: sessionService,
	}
}
func (s *Service) Execute(cmd interface{}) *models.Answer {
	switch cmd.(type) {
	case models.CommandCreateCredentials:
		err := s.credsService.InsertCredentials(cmd.(models.CommandCreateCredentials))
		return &models.Answer{
			Err: err,
		}
	case models.QueryReadCredentialsByUsername:
		credentials, err := s.credsService.GetCredentials(cmd.(models.QueryReadCredentialsByUsername))
		return &models.Answer{
			Err:         err,
			Credentials: credentials,
		}
	case models.CommandDeleteCredentialsByUsername:
		err := s.credsService.DeleteCredentials(cmd.(models.CommandDeleteCredentialsByUsername))
		return &models.Answer{
			Err: err,
		}
	case models.CommandCreateSession:
		err := s.sessionService.InsertSession(cmd.(models.CommandCreateSession))
		return &models.Answer{
			Err: err,
		}
	case models.QueryReadSessionByRefreshToken:
		session, err := s.sessionService.GetSession(cmd.(models.QueryReadSessionByRefreshToken))
		return &models.Answer{
			Err:     err,
			Session: session,
		}
	case models.CommandDeleteSessionByRefreshToken:
		err := s.sessionService.DeleteSession(cmd.(models.CommandDeleteSessionByRefreshToken))
		return &models.Answer{
			Err: err,
		}
	default:
		return &models.Answer{Err: ErrComandNotFound}
	}
}
