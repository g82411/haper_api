package paths

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
)

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

func resolveStyle(style int) string {
	styles := []string{
		"卡通插畫",
		"單一線條、色塊",
		"平面插畫",
		"擬真",
	}
	return styles[style]
}

func resolveAction(action int) string {
	actions := []string{
		"",
		"物品產生",
		"人物產生",
		"物品間關係",
		"台灣小吃",
		"節慶用品",
		"運動",
	}
	return actions[action]
}

func GenerateImage(c *fiber.Ctx) error {
	const PromptTemplate = "%v以卡通插圖的風格繪製，線條乾淨俐落，線條較粗，避免複雜、多餘的線條，使用簡單的色彩。\n圖片僅有主體、呈現完整的樣貌、貼近實際場景、貼近實際動作。主體為彩色，圖片背景是白色。"
	accessData := c.Locals("accessData").(utils.Claims)
	var body GenerateRequest
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	dbClient, err := models.NewDBClient()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	userSub := accessData.Sub

	hasEnoughCredit, err := bussinessLogic.CheckUserHaveEnoughCredit(userSub)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	if !hasEnoughCredit {
		c.Status(fiber.StatusPaymentRequired)
		return c.JSON(fiber.Map{
			"message": "今日額度已用完",
		})
	}
	prompt := fmt.Sprintf(PromptTemplate, body.Prompt)
	fmt.Printf("Prompt: %v\n", prompt)

	articleRecord := models.Article{
		ID:         utils.GenerateShortKey(),
		Tool:       "卡通插圖易讀產生",
		Style:      "卡通插畫",
		Keyword:    body.Prompt,
		AuthorId:   accessData.Sub,
		Valid:      false,
		AuthorName: accessData.Name,
	}
	taskRecord := models.Task{
		ID:       utils.GenerateShortKey(),
		Prompt:   prompt,
		AuthorId: accessData.Sub,
		Status:   0,
	}
	err = dbClient.Create(&articleRecord).Error
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	err = dbClient.Create(&taskRecord).Error
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	err = bussinessLogic.PutImageRequestToQueue(
		taskRecord.ID,
		accessData.Sub,
		prompt,
		articleRecord.ID,
	)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	needSurvey, err := bussinessLogic.CheckNeedJumpSurveyToUser(userSub)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"article_id":  articleRecord.ID,
		"task_id":     taskRecord.ID,
		"need_survey": needSurvey,
	})
}
