package bussinessLogic

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

func UpdateArticle(ctx context.Context, articleId, userId, ta, willUse string) error {
	updateQuery := "set ta = :ta, will_use = :will_use"
	key := map[string]types.AttributeValue{
		"date_id":   &types.AttributeValueMemberS{Value: articleId},
		"author_id": &types.AttributeValueMemberS{Value: userId},
	}
	expressionAttribute := map[string]types.AttributeValue{
		":ta":       &types.AttributeValueMemberS{Value: ta},
		":will_use": &types.AttributeValueMemberS{Value: willUse},
	}
	var article models.Article
	tableName := article.TableName(ctx)
	err := dynamodb.Update(ctx, tableName, &dynamodb.UpdateQuery{
		ExpressionAttribute: &expressionAttribute,
		UpdateExpression:    &updateQuery,
		Key:                 &key,
	})
	if err != nil {
		return err
	}
	return nil
}
