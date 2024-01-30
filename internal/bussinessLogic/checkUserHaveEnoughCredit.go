package bussinessLogic

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
	"hyper_api/internal/utils/resolver/date"
)

const DailyCredit = 10

func CheckUserHaveEnoughCredit(ctx context.Context, authorId string) (bool, error) {
	tableName := models.Article{}.TableName(ctx)
	todayDate := date.GetNowDateString()

	keyConditionExpression := "author_id = :author_id AND begins_with(date_id, :date)"
	expressionAttrVals := map[string]types.AttributeValue{
		":author_id": &types.AttributeValueMemberS{Value: authorId},
		":date":      &types.AttributeValueMemberS{Value: todayDate},
		":valid":     &types.AttributeValueMemberS{Value: "true"},
	}
	filter := "valid = :valid"
	query := dynamodb.InputQuery{
		KeyConditionExpression: &keyConditionExpression,
		ExpressionAttribute:    &expressionAttrVals,
		FilterExpression:       &filter,
	}
	items, err := dynamodb.Query(ctx, tableName, &query)
	if err != nil {
		return false, fmt.Errorf("error while query article: %s", err.Error())
	}
	n := len(items)
	return n <= DailyCredit, nil
}
