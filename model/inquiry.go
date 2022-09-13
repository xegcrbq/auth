package model

type RepositoryInquiry struct {
	Data                  []byte //возможно есть смысл поменять на interface
	RepositoryInquiryCode RepositoryInquiryCode
}

type RepositoryInquiryCode byte

const (
	CREATE RepositoryInquiryCode = iota
	READ
	UPDATE
	DELETE
)

type ServiceInquiry struct {
	ServiceInquiryCode ServiceInquiryCode
	SubInquiry         RepositoryInquiry
}
type ServiceInquiryCode byte

const (
	REFRESHSESSION ServiceInquiryCode = iota
	CREDENTIALS
)
