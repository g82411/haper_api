package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/models"
)

func GetArticle(c *fiber.Ctx) error {
	db, err := models.NewDBClient()
	err = db.AutoMigrate(&models.Article{})
	//var article models.Article
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	articleId := c.Params("articleId")
	var article models.Article
	tx := db.Preload("Tags").Table("articles").Where("id = ?", articleId).Find(&article)
	if tx.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return tx.Error
	}
	if article.ID == "" {
		res := make(map[string]interface{})
		res["message"] = "Article not found"
		c.Status(fiber.StatusNotFound)
		return c.JSON(res)
	}
	return c.JSON(article)
}
