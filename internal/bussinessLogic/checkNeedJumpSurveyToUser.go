package bussinessLogic

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

func CheckNeedJumpSurveyToUser(ctx context.Context, authorId string) (bool, error) {
	ArticleCountThatNeedSurvey := map[int]bool{
		5:  true,
		20: true,
		40: true,
	}
	const MaxQueryLimit = int32(41)
	tableName := models.Article{}.TableName(ctx)

	keyConditionExpression := "author_id = :author_id"
	expressionAttrVals := map[string]types.AttributeValue{
		":author_id": &types.AttributeValueMemberS{Value: authorId},
		":valid":     &types.AttributeValueMemberS{Value: "true"},
	}
	filter := "valid = :valid"
	query := dynamodb.InputQuery{
		KeyConditionExpression: &keyConditionExpression,
		ExpressionAttribute:    &expressionAttrVals,
		FilterExpression:       &filter,
		Limit:                  MaxQueryLimit,
	}
	items, err := dynamodb.Query(ctx, tableName, &query)
	if err != nil {
		return false, fmt.Errorf("error while query article: %s", err.Error())
	}
	n := len(items)

	return ArticleCountThatNeedSurvey[n], nil
}
