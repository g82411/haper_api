package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/dto"
)

func GetUserInfo(c *fiber.Ctx) error {
	userInfo := c.Locals("userInfo").(*dto.UserInfo)
	return c.JSON(userInfo)
}
