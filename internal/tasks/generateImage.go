package tasks

import (
	"context"
	"fmt"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws/dynamodb"
)

func GenerateImageToOpenAI(task dto.GenerateImageTask) error {
	imageUrl, err := bussinessLogic.GenerateImageByPrompt(task.Prompt)
	if err != nil {
		return fmt.Errorf("generate Image error %v", err)
	}
	ctx := context.Background()
	stageCtx := context.WithValue(ctx, "stage", "prod")
	dynamoCtx, err := dynamodb.WithDynamoDBConnection(stageCtx)
	if err != nil {
		return fmt.Errorf("error when connect to dynamodb %v", err)
	}
	err = bussinessLogic.SaveImageToArticle(dynamoCtx, task, imageUrl)
	if err != nil {
		return fmt.Errorf("error when save image to article %v", err)
	}
	return nil
}
