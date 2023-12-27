package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/models"
	"strconv"
)

const PerPage = 20

func TakeImages(c *fiber.Ctx) error {
	page := c.Query("page")
	pageInt := 1
	if page != "" {
		pageParsed, err := strconv.Atoi(page)
		if err == nil {
			pageInt = pageParsed
		}
	}
	db, err := models.NewDBClient()
	//var article models.Article
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	var articles []models.Article
	tx := db.Select("id", "url", "author_name", "keyword")
	tx.Table("articles")
	tx.Order("created_at desc")
	tx.Where("valid = ?", true)
	tx.Offset((pageInt - 1) * PerPage).Limit(PerPage).Find(&articles)
	result := []map[string]interface{}{}
	for _, article := range articles {
		result = append(result, map[string]interface{}{
			"id":          article.ID,
			"url":         article.Url,
			"author_name": article.AuthorName,
			"keyword":     article.Keyword,
		})
	}
	return c.JSON(fiber.Map{
		"articles": result,
	})
}
