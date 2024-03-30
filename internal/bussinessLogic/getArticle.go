package bussinessLogic

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

func GetArticle(ctx context.Context, articleId string) (map[string]interface{}, error) {
	articleKeyExpression := "id = :id"
	var articleFake models.Article
	articleQuery := dynamodb.InputQuery{
		KeyConditionExpression: &articleKeyExpression,
		ExpressionAttribute: &map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: articleId},
		},
		Limit:     1,
		IndexName: "PrimaryIdIndex",
	}
	article, err := dynamodb.Query(ctx, articleFake.TableName(ctx), &articleQuery)
	if err != nil {
		return nil, err
	}
	if len(article) == 0 {
		return nil, nil
	}
	articleT := articleFake.Serialize(article[0]).(models.Article)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          articleT.ID,
		"imageUrl":    articleT.Url,
		"keyword":     articleT.Keyword,
		"authorId":    articleT.AuthorId,
		"authorName":  articleT.AuthorName,
		"dateId":      articleT.DateId,
		"date":        articleT.Date,
		"tags":        articleT.Age,
		"region":      articleT.Region,
		"ta":          articleT.TA,
		"willUse":     articleT.WillUse,
		"authorImage": articleT.AuthorImage,
	}, nil
}
