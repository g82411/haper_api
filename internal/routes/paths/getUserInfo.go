package paths

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws"
)

func GetUserInfo(c *fiber.Ctx) error {
	previousInfo := c.Locals("accessData").(utils.Claims)
	accessToken := c.Locals("accessToken").(string)
	userSub := previousInfo.Sub
	svc, err := aws.NewDynamoDBClient()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	_, refreshToken, err := aws.GetTokenFromDB(svc, userSub)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	renewToken, err := utils.RefreshToken(accessToken, refreshToken, true)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return err
	}
	userInfo, err := utils.ExtractUserInfoFromToken(renewToken.Extra("id_token").(string))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	return c.JSON(userInfo)
}
