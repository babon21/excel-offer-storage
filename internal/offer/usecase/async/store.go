package usecase

type Store interface {
	Get(id int64) (string, error)
	Set(id int64, value string) error
	GetNewId(idFieldName string) (int64, error)
}
