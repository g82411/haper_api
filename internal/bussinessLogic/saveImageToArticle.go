package bussinessLogic

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/dto"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

func SaveImageToArticle(dynamoCtx context.Context, task dto.GenerateImageTask, imageUrl string) error {
	updateQuery := "set image_url = :imageUrl"
	key := map[string]types.AttributeValue{
		"date_id":   &types.AttributeValueMemberS{Value: task.ArticleDateId},
		"author_id": &types.AttributeValueMemberS{Value: task.AuthorId},
	}
	expressionAttribute := map[string]types.AttributeValue{
		":imageUrl": &types.AttributeValueMemberS{Value: imageUrl},
	}
	var article models.Article
	tableName := article.TableName(dynamoCtx)
	err := dynamodb.Update(dynamoCtx, tableName, &dynamodb.UpdateQuery{
		ExpressionAttribute: &expressionAttribute,
		UpdateExpression:    &updateQuery,
		Key:                 &key,
	})
	if err != nil {
		return err
	}
	return nil
}
