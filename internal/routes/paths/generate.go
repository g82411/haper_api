package paths

import (
	context2 "context"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws/dynamodb"
)

type GenerateRequest struct {
	Prompt string   `json:"prompt"`
	Tags   []string `json:"tags"`
}

func GenerateImage(c *fiber.Ctx) error {
	userInfo := c.Locals("userInfo").(*dto.UserInfo)
	//isLocal := c.Locals("env:isLocal").(bool)
	var body GenerateRequest
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	ctx := context2.Background()
	dynamoCtx, err := dynamodb.WithDynamoDBConnection(ctx)
	stageCtx := context2.WithValue(dynamoCtx, "stage", "prod")
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	hasCredit, err := bussinessLogic.CheckUserHaveEnoughCredit(stageCtx, userInfo.Sub)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	if !hasCredit {
		c.Status(fiber.StatusPaymentRequired)
		return c.JSON(fiber.Map{
			"message": "今日額度已用完",
		})
	}
	article, err := bussinessLogic.CreateArticle(stageCtx, userInfo, body.Prompt, body.Tags)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	err = bussinessLogic.InsertRelationBetweenTagAndArticle(stageCtx, article, body.Tags)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	needJump := false
	needJump, err = bussinessLogic.CheckNeedJumpSurveyToUser(stageCtx, userInfo.Sub)
	if err != nil {
		needJump = false
	}

	task, err := bussinessLogic.CreateTask(stageCtx, article)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	err = bussinessLogic.PutImageRequestToQueue(task)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"article_id":  article.ID,
		"task_id":     task.ID,
		"need_survey": needJump,
	})
}
