package paths

import (
	context2 "context"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/utils/aws/dynamodb"
)

func GetArticle(c *fiber.Ctx) error {
	articleId := c.Params("articleId")
	ctx := context2.Background()
	dynamoCtx, err := dynamodb.WithDynamoDBConnection(ctx)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	stageCtx := context2.WithValue(dynamoCtx, "stage", "prod")
	article, err := bussinessLogic.GetArticle(stageCtx, articleId)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	if article == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Article not found",
		})
	}
	return c.JSON(article)
}
