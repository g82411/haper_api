package models

type Tag struct {
	ID      string     `gorm:"type:uuid;primary_key;"`
	Name    string     `gorm:"unique"`
	Article []*Article `gorm:"many2many:article_tags;"`
}
