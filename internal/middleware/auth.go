package middleware

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
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
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	svc := dynamodb.NewFromConfig(cfg)
	accessData, err := aws.GetUserInfoByAccessToken(svc, token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	c.Locals("accessData", accessData)
	c.Locals("accessToken", token)
	return c.Next()
}
