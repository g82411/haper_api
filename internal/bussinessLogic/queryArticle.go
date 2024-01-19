package bussinessLogic

import (
	"context"
	"gorm.io/gorm"
	"hyper_api/internal/models"
)

type scope func(*gorm.DB) *gorm.DB

func QueryArticle(ctx context.Context, page, pageSize int, fn scope) ([]map[string]interface{}, error) {
	db := ctx.Value("db").(*gorm.DB)
	var articles []models.Article
	tx := db.Preload("Tags").Select("id", "url", "author_name", "keyword")
	tx.Table("articles")
	tx.Order("created_at desc")
	tx.Where("valid = ?", true)
	tx = fn(tx)
	tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)
	result := []map[string]interface{}{}
	for _, article := range articles {
		tagEntities := article.Tags
		tags := make([]string, 0)
		for _, tagEntity := range tagEntities {
			tags = append(tags, tagEntity.Name)
		}
		result = append(result, map[string]interface{}{
			"id":          article.ID,
			"url":         article.Url,
			"author_name": article.AuthorName,
			"keyword":     article.Keyword,
			"tags":        tags,
		})
	}
	return result, nil
}
