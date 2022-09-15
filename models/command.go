package models

type Command interface {
	string
	IsCorrect() bool
}

// session
type QueryReadSessionByRefreshToken struct {
	RefreshToken string
}

func (c *QueryReadSessionByRefreshToken) IsCorrect() bool {
	return c.RefreshToken != ""
}

type CommandCreateSession struct {
	Session *Session
}

func (c *CommandCreateSession) IsCorrect() bool {
	return c.Session.IsValid()
}

type CommandDeleteSessionByRefreshToken struct {
	RefreshToken string
}

func (c *CommandDeleteSessionByRefreshToken) IsCorrect() bool {
	return c.RefreshToken != ""
}

// credentials
type QueryReadCredentialsByUsername struct {
	Username string
}

func (c *QueryReadCredentialsByUsername) IsCorrect() bool {
	return c.Username != ""
}

type CommandCreateCredentials struct {
	Credentials *Credentials
}

func (c *CommandCreateCredentials) IsCorrect() bool {
	return c.Credentials.IsValid()
}

type CommandDeleteCredentialsByUsername struct {
	Username string
}

func (c *CommandDeleteCredentialsByUsername) IsCorrect() bool {
	return c.Username != ""
}
