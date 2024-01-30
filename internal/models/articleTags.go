package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/utils/aws/dynamodb"
)

type ArticleTag struct {
	dynamodb.SerializeAble
	ArticleID string
	TagName   string
}

func (articleTag ArticleTag) TableName(ctx context.Context) string {
	stage := ctx.Value("stage").(string)
	return stage + "_article_tags"
}

func (articleTag ArticleTag) Deserialize() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"article_id": &types.AttributeValueMemberS{Value: articleTag.ArticleID},
		"tag_name":   &types.AttributeValueMemberS{Value: articleTag.TagName},
	}
}

func (_ ArticleTag) Serialize(av map[string]types.AttributeValue) interface{} {
	var articleTag ArticleTag
	articleTag.ArticleID = av["article_id"].(*types.AttributeValueMemberS).Value
	articleTag.TagName = av["tag_name"].(*types.AttributeValueMemberS).Value
	return articleTag
}
