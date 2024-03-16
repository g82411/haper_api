package models

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"hyper_api/internal/utils/aws/dynamodb"
)

type Article struct {
	dynamodb.SerializeAble
	ID          string
	Url         string
	Keyword     string
	AuthorId    string
	AuthorName  string
	Date        string
	Valid       bool
	DateId      string
	AuthorImage string
	Tags        string
}

func (article Article) TableName(ctx context.Context) string {
	stage := ctx.Value("stage").(string)
	return fmt.Sprintf("%s_articles", stage)
}

func (article Article) Deserialize() map[string]types.AttributeValue {
	valid := "false"
	if article.Valid {
		valid = "true"
	}
	return map[string]types.AttributeValue{
		"id":           &types.AttributeValueMemberS{Value: article.ID},
		"image_url":    &types.AttributeValueMemberS{Value: article.Url},
		"keyword":      &types.AttributeValueMemberS{Value: article.Keyword},
		"author_id":    &types.AttributeValueMemberS{Value: article.AuthorId},
		"author_name":  &types.AttributeValueMemberS{Value: article.AuthorName},
		"author_image": &types.AttributeValueMemberS{Value: article.AuthorImage},
		"tags":         &types.AttributeValueMemberS{Value: article.Tags},
		"date":         &types.AttributeValueMemberS{Value: article.Date},
		"valid":        &types.AttributeValueMemberS{Value: valid},
		"date_id":      &types.AttributeValueMemberS{Value: article.DateId},
	}
}

func (_ Article) Serialize(av map[string]types.AttributeValue) interface{} {
	var article Article
	article.ID = av["id"].(*types.AttributeValueMemberS).Value
	article.Url = av["image_url"].(*types.AttributeValueMemberS).Value
	article.Keyword = av["keyword"].(*types.AttributeValueMemberS).Value
	article.AuthorId = av["author_id"].(*types.AttributeValueMemberS).Value
	article.AuthorName = av["author_name"].(*types.AttributeValueMemberS).Value
	article.Valid = av["valid"].(*types.AttributeValueMemberS).Value == "true"
	article.DateId = av["date_id"].(*types.AttributeValueMemberS).Value
	return article
}
