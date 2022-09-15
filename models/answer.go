package models

type Answer struct {
	Err         error
	Session     *Session
	Credentials *Credentials
}
