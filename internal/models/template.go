package models

type Template struct {
	ID        string `gorm:"primaryKey"`
	Prompt    string
	Title     string
	DeletedAt int64
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}
