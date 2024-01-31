package bussinessLogic

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

const PerPage = 30

type QueryOption struct {
	AuthorId     string
	LastAuthorId string
	LastDateId   string
}

func transformOptToQuery(opt *QueryOption) *dynamodb.InputQuery {
	var query dynamodb.InputQuery
	keyCondition := "valid = :valid"
	expressionAttrVals := map[string]types.AttributeValue{
		":valid": &types.AttributeValueMemberS{Value: "true"},
	}
	query.IndexName = "ValidDateIndex"
	query.ScanIndexForward = false
	query.Limit = PerPage
	if opt.LastDateId != "" && opt.LastAuthorId != "" {
		query.ExclusiveStartKey = &map[string]types.AttributeValue{
			"valid":     &types.AttributeValueMemberS{Value: "true"},
			"author_id": &types.AttributeValueMemberS{Value: opt.LastAuthorId},
			"date_id":   &types.AttributeValueMemberS{Value: opt.LastDateId},
		}
	}
	if opt.AuthorId != "" {
		query.IndexName = ""
		keyCondition = "author_id = :author_id"
		expressionAttrVals[":author_id"] = &types.AttributeValueMemberS{Value: opt.AuthorId}
		filterExpression := "valid = :valid"
		query.FilterExpression = &filterExpression
		delete(expressionAttrVals, ":valid")
	}
	query.KeyConditionExpression = &keyCondition
	query.ExpressionAttribute = &expressionAttrVals
	return &query
}

func QueryArticle(ctx context.Context, opt *QueryOption) ([]map[string]interface{}, error) {
	var article models.Article
	tableName := article.TableName(ctx)
	query := transformOptToQuery(opt)
	items, err := dynamodb.Query(ctx, tableName, query)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	for _, v := range items {
		id, _ := v["id"].(*types.AttributeValueMemberS)
		authorId, _ := v["author_id"].(*types.AttributeValueMemberS)
		dateId, _ := v["date_id"].(*types.AttributeValueMemberS)
		authorName, _ := v["author_name"].(*types.AttributeValueMemberS)
		keyword, _ := v["keyword"].(*types.AttributeValueMemberS)
		imageUrl, _ := v["image_url"].(*types.AttributeValueMemberS)
		result = append(result, map[string]interface{}{
			"id":         id.Value,
			"authorId":   authorId.Value,
			"dateId":     dateId.Value,
			"authorName": authorName.Value,
			"keyword":    keyword.Value,
			"imageUrl":   imageUrl.Value,
		})
	}
	return result, nil
}
