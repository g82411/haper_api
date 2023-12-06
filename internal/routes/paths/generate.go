package paths

import (
	"context"
	"fmt"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
	Relation string   `json:"relation"`
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
	token := c.Cookies("token")
	setting := config.GetConfig()
	fmt.Println("token", token)
	if token == "" {
		return c.Redirect(config.GetConfig().CognitoDomain + "/oauth2/authorize?response_type=code&client_id=" + config.GetConfig().CognitoClientId + "&redirect_uri=" + config.GetConfig().RedirectURL)
	}
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	svc := dynamodb.NewFromConfig(cfg)
	accessData, err := aws.GetUserInfoByAccessToken(svc, token)
	var body GenerateRequest
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	var uploadReq []aws.PutObjectInput
	prompt := resolver.GenerateImagePrompt(dto.GenerateImageRequest{
		Items:    body.Items,
		Relation: body.Relation,
		Action:   body.Action,
		Style:    body.Style,
	})
	//var generateImageUrls []string
	generateImageUrls := make([]string, 1)

	ans, _ := utils.GeneratePhotoUsingDallE3(prompt, 1)
	if ans == nil || len(ans) == 0 {
		return c.JSON(fiber.Map{
			"Error": "Generate image failed",
		})
	}
	generateImageUrls[0] = ans[0]
	downloadResult, _ := utils.DownloadImages(generateImageUrls)
	for _, v := range downloadResult {
		if !v.Success {
			log.Error("Download image failed: ", v.Error)
		}
		uploadReq = append(uploadReq, aws.PutObjectInput{
			Bucket:      setting.GenerateS3Bucket,
			Key:         fmt.Sprintf("%v.png", uuid.New().String()),
			Body:        v.Image,
			ContentType: "image/png",
		})
	}
	sess, _ := aws.NewAWSSession()
	s3Client := aws.NewS3Client(sess)
	uploadedImages, err := s3Client.PutObjects(uploadReq)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	var outputUrls []string
	for _, v := range uploadedImages {
		if !v.Success {
			log.Error("Upload image failed: ", v.Error)
		}
		outputUrls = append(outputUrls, fmt.Sprintf("https://%v/%v", setting.CDNHost, v.Key))
	}

	outputUrl := outputUrls[0]
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
		Url:        outputUrl,
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
