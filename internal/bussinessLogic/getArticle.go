package bussinessLogic

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

func GetArticle(ctx context.Context, articleId string) (map[string]interface{}, error) {
	articleKeyExpression := "id = :id"
	articleFilterExpression := "valid = :valid"
	var articleFake models.Article
	var tags models.ArticleTag
	articleQuery := dynamodb.InputQuery{
		KeyConditionExpression: &articleKeyExpression,
		ExpressionAttribute: &map[string]types.AttributeValue{
			":id":    &types.AttributeValueMemberS{Value: articleId},
			":valid": &types.AttributeValueMemberS{Value: "true"},
		},
		FilterExpression: &articleFilterExpression,
		Limit:            1,
		IndexName:        "PrimaryIdIndex",
	}
	article, err := dynamodb.Query(ctx, articleFake.TableName(ctx), &articleQuery)
	if err != nil {
		return nil, err
	}
	if len(article) == 0 {
		return nil, nil
	}
	articleT := articleFake.Serialize(article[0]).(models.Article)

	tagKeyExpression := "article_id = :id"
	tagsQuery := dynamodb.InputQuery{
		KeyConditionExpression: &tagKeyExpression,
		ExpressionAttribute: &map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: articleId},
		},
		IndexName: "ArticleIdIndex",
	}
	tagsResult, err := dynamodb.Query(ctx, tags.TableName(ctx), &tagsQuery)
	if err != nil {
		return nil, err
	}
	tag := make([]string, len(tagsResult))
	for i, v := range tagsResult {
		tag[i] = v["tag_name"].(*types.AttributeValueMemberS).Value
	}
	return map[string]interface{}{
		"id":         articleT.ID,
		"imageUrl":   articleT.Url,
		"keyword":    articleT.Keyword,
		"authorId":   articleT.AuthorId,
		"authorName": articleT.AuthorName,
		"dateId":     articleT.DateId,
		"date":       articleT.Date,
		"tags":       tag,
	}, nil
}
