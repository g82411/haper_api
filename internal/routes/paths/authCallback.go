package paths

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws"
)

func AuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := utils.ExtractTokenFromCode(code)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	idToken, _ := token.Extra("id_token").(string)
	userData, err := utils.ExtractUserInfoFromToken(idToken)
	if userData.IsDoneSurvey != "true" {
		err = aws.UpdateUserSurveyStatus(userData.CogUsername, false)
	}
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	svc := dynamodb.NewFromConfig(cfg)
	err = aws.PutTokenToDB(svc, userData.Sub, idToken, token.RefreshToken)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token.AccessToken
	cookie.HTTPOnly = true
	c.Cookie(cookie)
	res := map[string]interface{}{
		"username": userData.Name,
		"email":    userData.Email,
	}
	c.Status(fiber.StatusCreated)
	if userData.IsDoneSurvey == "true" {
		c.Status(fiber.StatusOK)
	}
	return c.JSON(res)
}
