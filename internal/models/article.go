package models

type Article struct {
	ID         string `gorm:"primaryKey"`
	Url        string
	Tool       string
	Style      string
	Keyword    string
	AuthorId   string
	AuthorName string
	Valid      bool `gorm:"default:true"`
	DeletedAt  int64
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
	Tags       []Tag `gorm:"many2many:article_tags;"`
}
