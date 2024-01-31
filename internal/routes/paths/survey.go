package paths

import (
	"github.com/gofiber/fiber/v2"
)

type SurveyPayload struct {
	Reason string `json:"reason"`
	Age    string `json:"age"`
	Career string `json:"career"`
	Gender string `json:"gender"`
	Name   string `json:"name"`
}

func Survey(c *fiber.Ctx) error {
	//userInfo := c.Locals("userInfo").(*dto.UserInfo)
	//cognitoPoolId := c.Locals("config:CognitoUserPoolId").(string)
	var body SurveyPayload
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	//svc, err := aws.NewDynamoDBClient()
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//err = aws.PutSurveyResultToDB(svc, userInfo.Sub, body.Reason)
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	//client, err := aws.NewAdminCognitoClient(cognitoPoolId)
	//err = client.UpdateUserInfo(
	//	userInfo.InternalUserName,
	//	body.Name,
	//	body.Age,
	//	body.Career,
	//	body.Gender,
	//)
	//if err != nil {
	//	c.Status(fiber.StatusInternalServerError)
	//	return err
	//}
	return c.SendStatus(fiber.StatusCreated)
}
