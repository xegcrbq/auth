package auth

type Repository interface {
	AddData(data dbData) (dbData, error)
	RemoveData(data dbData) (dbData, error)
	GetExistingData(data dbData) (dbData, error)
}
