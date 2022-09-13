package model

// Credentials структура для парсинга данных из json и бд
type Credentials struct {
	UserId   int32  `json:"userId" db:"userId"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"userName"`
}

func (c Credentials) IsValid() bool {
	if c.Username != "" && c.Password != "" {
		return true
	}
	return false
}
