package paths

import (
	context2 "context"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws/dynamodb"
)

type UpdateRequest struct {
	TA      string `json:"ta"`
	WillUse string `json:"willUse"`
}

func UpdateArticle(c *fiber.Ctx) error {
	articleId := c.Params("articleId")
	var body UpdateRequest
	userInfo := c.Locals("userInfo").(*dto.UserInfo)

	if err := c.BodyParser(&body); err != nil {
		return err
	}
	ctx := context2.Background()
	dynamoCtx, err := dynamodb.WithDynamoDBConnection(ctx)
	stageCtx := context2.WithValue(dynamoCtx, "stage", "prod")
	err = bussinessLogic.UpdateArticle(stageCtx, articleId, userInfo.Sub, body.TA, body.WillUse)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{})
}
