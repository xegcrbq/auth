package model

type Answer struct {
	Data []byte
	Code AnswerCode
}

type AnswerCode byte

const (
	SUCCSESS AnswerCode = iota
	DATACREDENTIALS
	DATAREFRESHSESSION
	ERROR
)
