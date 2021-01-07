package db

type Model struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
