package models

// Credentials структура для данных пользователя
type Credentials struct {
	UserId   int32  `json:"userid" db:"userid"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
}

func (c Credentials) IsValid() bool {
	if c.Username != "" && c.Password != "" {
		return true
	}
	return false
}
func (c Credentials) Equal(c2 Credentials) bool {
	if c.Username != c2.Username {
		return false
	}
	if c.Password != c2.Password {
		return false
	}
	return true
}

type CredentialsRepository interface {
	SaveCredentials(c *Credentials) error
	ReadCredentialsByUsername(username string) (*Credentials, error)
	DeleteCredentialsByUsername(username string) error
}
