package paths

import (
	context2 "context"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/utils/aws/dynamodb"
)

func TakeImages(c *fiber.Ctx) error {
	lastAuthorId := c.Query("last_author_id")
	lastDateId := c.Query("last_date_id")

	ctx := context2.Background()
	dynamoCtx, err := dynamodb.WithDynamoDBConnection(ctx)
	stageCtx := context2.WithValue(dynamoCtx, "stage", "prod")
	result, err := bussinessLogic.QueryArticle(stageCtx, &bussinessLogic.QueryOption{
		LastAuthorId: lastAuthorId,
		LastDateId:   lastDateId,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{
		"articles": result,
	})
}
