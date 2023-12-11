package paths

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws"
	"hyper_api/internal/utils/resolver"
)

type GenerateRequest struct {
	Items    []string `json:"items"`
	Relation string   `json:"relationship"`
	Action   int      `json:"action"`
	Style    int      `json:"style"`
}

func resolveStyle(style int) string {
	styles := []string{
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
	}
	return actions[action]
}

func GenerateImage(c *fiber.Ctx) error {
	setting := config.GetConfig()
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
	})
	//var generateImageUrls []string
	fmt.Printf("Prompt: %v\n", prompt)
	ans, _ := utils.GeneratePhotoUsingDallE3(prompt, 1)
	if ans == nil || len(ans) == 0 {
		return c.JSON(fiber.Map{
			"Error": "Generate image failed",
		})
	}
	url := ans[0]
	downloadResult := utils.DownloadImage(url)
	uploadReq := aws.PutObjectInput{
		Bucket:      setting.GenerateS3Bucket,
		Key:         fmt.Sprintf("%v.png", uuid.New().String()),
		Body:        downloadResult.Image,
		ContentType: "image/png",
	}
	sess, _ := aws.NewAWSSession()
	s3Client := aws.NewS3Client(sess)
	uploadedImages := s3Client.PutObject(uploadReq)
	if uploadedImages.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return uploadedImages.Error
	}
	newImageUrl := fmt.Sprintf("%v/%v", setting.CDNHost, uploadedImages.Key)
	keyword := body.Items[0]
	if body.Action == 3 {
		keyword = fmt.Sprintf("%v在%v的%v", body.Items[0], body.Items[1], body.Relation)
	}
	dbClient, err := models.NewDBClient()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	// TODO: Remove this line before deploy
	err = dbClient.AutoMigrate(&models.Article{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	articleRecord := models.Article{
		ID:         utils.GenerateShortKey(),
		Url:        newImageUrl,
		Tool:       resolveAction(body.Action),
		Style:      resolveStyle(body.Style),
		Keyword:    keyword,
		AuthorId:   accessData.Sub,
		AuthorName: accessData.Name,
	}
	err = dbClient.Create(&articleRecord).Error
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{
		"id":  articleRecord.ID,
		"url": articleRecord.Url,
	})
}
