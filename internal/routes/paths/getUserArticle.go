package paths

import (
	context2 "context"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws/dynamodb"
)

func QueryArticleByUser(c *fiber.Ctx) error {
	page := c.Query("last_date_id")
	userInfo := c.Locals("userInfo").(*dto.UserInfo)
	userSub := userInfo.Sub
	ctx := context2.Background()
	dynamoCtx, err := dynamodb.WithDynamoDBConnection(ctx)
	stageCtx := context2.WithValue(dynamoCtx, "stage", "prod")
	result, err := bussinessLogic.QueryArticle(stageCtx, &bussinessLogic.QueryOption{
		LastAuthorId: userSub,
		LastDateId:   page,
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
