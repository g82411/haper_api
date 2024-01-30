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
	const PromptTemplate = "%v以卡通插圖的風格繪製，線條乾淨俐落，線條較粗，避免複雜、多餘的線條，使用簡單的色彩。\n圖片僅有主體、呈現完整的樣貌、貼近實際場景、貼近實際動作。主體為彩色，圖片背景是白色。"
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
	article, err := bussinessLogic.CreateArticle(stageCtx, userInfo.Sub, userInfo.Name, body.Prompt)
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

	//dbClient, err := models.NewDBClient()
	//tags := body.Tags
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//hasEnoughCredit, err := bussinessLogic.CheckUserHaveEnoughCredit(userInfo.Sub)
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//if !hasEnoughCredit {
	//	c.Status(fiber.StatusPaymentRequired)
	//	return c.JSON(fiber.Map{
	//		"message": "今日額度已用完",
	//	})
	//}
	//prompt := fmt.Sprintf(PromptTemplate, body.Prompt)
	//context := context2.Background()
	//context = models.NewDBClientWithContext(context)
	//tagRecords, err := bussinessLogic.FindOrCreateTags(context, tags)
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//articleRecord := models.Article{
	//	ID:         utils.GenerateShortKey(),
	//	Keyword:    body.Prompt,
	//	AuthorId:   userInfo.Sub,
	//	Valid:      false,
	//	AuthorName: userInfo.Name,
	//}
	//taskRecord := models.Task{
	//	ID:       utils.GenerateShortKey(),
	//	Prompt:   prompt,
	//	AuthorId: userInfo.Sub,
	//	Status:   0,
	//}
	//err = dbClient.Create(&articleRecord).Error
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//
	//err = dbClient.Create(&taskRecord).Error
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//if !isLocal {
	//	err = bussinessLogic.PutImageRequestToQueue(
	//		taskRecord.ID,
	//		userInfo.Sub,
	//		prompt,
	//		articleRecord.ID,
	//	)
	//}
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//needSurvey, err := bussinessLogic.CheckNeedJumpSurveyToUser(userInfo.Sub)
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}

	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"article_id": article.ID,
		//"task_id":     taskRecord.ID,
		"need_survey": needJump,
	})
}
