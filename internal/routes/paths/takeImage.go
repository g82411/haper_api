package paths

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hyper_api/internal/bussinessLogic"
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
	ctx := context.Context(c.Context())
	ctx = models.NewDBClientWithContext(ctx)
	result, err := bussinessLogic.QueryArticle(ctx, pageInt, PerPage, func(tx *gorm.DB) *gorm.DB {
		return tx
	})
	//var article models.Article
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{
		"articles": result,
	})
}
