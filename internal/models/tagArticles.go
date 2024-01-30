package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/utils/aws/dynamodb"
)

type TagArticle struct {
	dynamodb.SerializeAble
	ArticleID string
	TagName   string
}

func (tagArticle TagArticle) TableName(ctx context.Context) string {
	stage := ctx.Value("stage").(string)
	return stage + "_tag_articles"
}

func (tagArticle TagArticle) Deserialize() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"article_id": &types.AttributeValueMemberS{Value: tagArticle.ArticleID},
		"tag_name":   &types.AttributeValueMemberS{Value: tagArticle.TagName},
	}
}

func (_ TagArticle) Serialize(av map[string]types.AttributeValue) interface{} {
	var tagArticle TagArticle
	tagArticle.ArticleID = av["article_id"].(*types.AttributeValueMemberS).Value
	tagArticle.TagName = av["tag_name"].(*types.AttributeValueMemberS).Value
	return tagArticle
}
