package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/models"
)

func GetTemplates(c *fiber.Ctx) error {
	db, err := models.NewDBClient()
	//var template models.Template
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	var templates []models.Template
	tx := db.Debug().Table("templates").Find(&templates)
	if tx.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return tx.Error
	}
	return c.JSON(templates)
}
