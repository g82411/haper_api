package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"hyper_api/internal/utils/aws"
	"strings"
)

func extractToken(c *fiber.Ctx) string {
	// 从Authorization头部获取值
	authHeader := c.Get("Authorization")

	// 检查值是否以"Bearer "开头
	if strings.HasPrefix(authHeader, "Bearer ") {
		// 提取并返回token部分
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 如果没有找到有效的token，则返回空字符串或错误
	return ""
}

func AuthMiddleware(c *fiber.Ctx) error {
	token := extractToken(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	client, err := aws.NewCognitoClient()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		log.Errorf("Error loading AWS config: %v", err)
	}
	userInfo, err := client.GetUserAttributeByAccessToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		log.Errorf("Error get user data: %v", err)
		return c.JSON(fiber.Map{
			"message": "Invalid token",
		})
	}
	// Create the Cognito Identity Provider Client
	c.Locals("userInfo", userInfo)
	return c.Next()
}
