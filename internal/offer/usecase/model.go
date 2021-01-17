package usecase

type Statistic struct {
	CreatedCount uint32 `json:"created_count"`
	UpdatedCount uint32 `json:"updated_count"`
	DeletedCount uint32 `json:"deleted_count"`
	ErrCount     uint32 `json:"err_count"`
}
