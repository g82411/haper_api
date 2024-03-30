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
	Region      string
	TA          string
	WillUse     string
	Age         string
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
		"region":       &types.AttributeValueMemberS{Value: article.Region},
		"age":          &types.AttributeValueMemberS{Value: article.Age},
		"tags":         &types.AttributeValueMemberS{Value: article.Tags},
		"date":         &types.AttributeValueMemberS{Value: article.Date},
		"valid":        &types.AttributeValueMemberS{Value: valid},
		"date_id":      &types.AttributeValueMemberS{Value: article.DateId},
	}
}

func (_ Article) Serialize(av map[string]types.AttributeValue) interface{} {
	var article Article
	article.ID = av["id"].(*types.AttributeValueMemberS).Value
	region := ""
	age := ""
	ta := ""
	willUse := ""
	if av["ta"] != nil {
		ta = av["ta"].(*types.AttributeValueMemberS).Value
	}
	if av["will_use"] != nil {
		willUse = av["will_use"].(*types.AttributeValueMemberS).Value
	}
	if av["region"] != nil {
		region = av["region"].(*types.AttributeValueMemberS).Value
	}
	if av["age"] != nil {
		age = av["age"].(*types.AttributeValueMemberS).Value
	}
	authorImage := ""
	if av["author_image"] != nil {
		authorImage = av["author_image"].(*types.AttributeValueMemberS).Value
	}
	article.Region = region
	article.Age = age
	article.TA = ta
	article.WillUse = willUse
	article.Url = av["image_url"].(*types.AttributeValueMemberS).Value
	article.Keyword = av["keyword"].(*types.AttributeValueMemberS).Value
	article.AuthorId = av["author_id"].(*types.AttributeValueMemberS).Value
	article.AuthorName = av["author_name"].(*types.AttributeValueMemberS).Value
	article.AuthorImage = authorImage
	article.Valid = av["valid"].(*types.AttributeValueMemberS).Value == "true"
	article.DateId = av["date_id"].(*types.AttributeValueMemberS).Value
	return article
}
