package paths

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/resolver"
)

type GenerateRequest struct {
	Items    []string `json:"items"`
	Relation string   `json:"relationship"`
	Action   int      `json:"action"`
	Style    int      `json:"style"`
	Comment  string   `json:"comment"`
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
	accessData := c.Locals("accessData").(utils.Claims)
	var body GenerateRequest
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	prompt := resolver.GenerateImagePrompt(dto.GenerateImageRequest{
		Items:    body.Items,
		Relation: body.Relation,
		Action:   body.Action,
		Style:    body.Style,
		Comment:  body.Comment,
	})
	//var generateImageUrls []string
	fmt.Printf("Prompt: %v\n", prompt)

	keyword := body.Items[0]
	if body.Action == 3 {
		keyword = fmt.Sprintf("%v在%v的%v", body.Items[0], body.Items[1], body.Relation)
	}
	if body.Action == 4 {
		sauceContainer, sauceName := body.Items[0], body.Items[1]
		sauce := body.Comment
		saucePos := body.Relation
		sauceString := ""
		if sauce != "" && saucePos != "" {
			sauceString = fmt.Sprintf("，%v在%v", sauce, saucePos)
		}
		keyword = fmt.Sprintf("%v的%v%v", sauceContainer, sauceName, sauceString)
	}
	if body.Action == 5 {
		color, itemName, shape := body.Items[0], body.Items[1], body.Items[2]
		comment := body.Comment
		keyword = fmt.Sprintf("%v的%v，形狀為%v，%v", color, itemName, shape, comment)
	}
	if body.Action == 6 {
		sportName := body.Items[0]
		age := body.Comment
		if age != "" {
			age = fmt.Sprintf("，%v時期", age)
		}
		keyword = fmt.Sprintf("%v%v", sportName, age)
	}
	if body.Comment != "" {
		keyword = fmt.Sprintf("%v,%v", keyword, body.Comment)
	}
	dbClient, err := models.NewDBClient()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	articleRecord := models.Article{
		ID:         utils.GenerateShortKey(),
		Tool:       resolveAction(body.Action),
		Style:      resolveStyle(body.Style),
		Keyword:    keyword,
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
	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"article_id": articleRecord.ID,
		"task_id":    taskRecord.ID,
	})
}
