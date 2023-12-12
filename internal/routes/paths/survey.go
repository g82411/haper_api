package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws"
)

type SurveyPayload struct {
	Reason string `json:"reason"`
	Age    string `json:"age"`
	Career string `json:"career"`
	Gender string `json:"gender"`
	Name   string `json:"name"`
}

func Survey(c *fiber.Ctx) error {
	userSub := c.Locals("accessData").(utils.Claims).Sub
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
	err = aws.UpdateUserInfo(
		userInfo.CogUsername,
		body.Name,
		body.Age,
		body.Career,
		body.Gender,
	)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}
