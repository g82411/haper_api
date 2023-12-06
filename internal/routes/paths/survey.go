package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/config"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws"
)

type SurveyPayload struct {
	Reason string `json:"reason"`
}

func Survey(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Redirect(config.GetConfig().CognitoDomain + "/oauth2/authorize?response_type=code&client_id=" + config.GetConfig().CognitoClientId + "&redirect_uri=" + config.GetConfig().RedirectURL)
	}
	accessData, err := utils.ExtractUserInfoFromToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return err
	}
	userSub := accessData.Sub
	var body SurveyPayload
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	svc, err := aws.NewDynamoDBClient()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	err = aws.PutSurveyResultToDB(svc, userSub, body.Reason)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	idToken, _, err := aws.GetTokenFromDB(svc, userSub)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	userInfo, err := utils.ExtractUserInfoFromToken(idToken)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	err = aws.UpdateUserSurveyStatus(userInfo.CogUsername, true)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}
