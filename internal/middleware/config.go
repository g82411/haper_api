package middleware

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func LoadConfig(c *fiber.Ctx) error {
	c.Locals("config:CognitoUserPoolId", os.Getenv("COGNITO_USER_POOL"))
	return c.Next()
}
