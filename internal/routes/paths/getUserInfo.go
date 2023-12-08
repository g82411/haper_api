package paths

import "github.com/gofiber/fiber/v2"

func GetUserInfo(c *fiber.Ctx) error {
	userInfo := c.Locals("accessData")
	return c.JSON(userInfo)
}
