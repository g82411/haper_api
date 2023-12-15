package models

type Task struct {
	ID        string `gorm:"primaryKey"`
	Prompt    string
	AuthorId  string
	Status    int `gorm:"default:0"`
	DeletedAt int64
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}
