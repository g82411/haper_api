package models

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Task struct {
	ID            string
	Prompt        string
	Age           string
	Region        string
	AuthorId      string
	ArticleDateId string
	Valid         bool
}

func (task *Task) TableName(ctx context.Context) string {
	stage := ctx.Value("stage").(string)
	return fmt.Sprintf("%s_tasks", stage)
}

func (task *Task) Deserialize() (map[string]types.AttributeValue, error) {
	return map[string]types.AttributeValue{
		"id":              &types.AttributeValueMemberS{Value: task.ID},
		"prompt":          &types.AttributeValueMemberS{Value: task.Prompt},
		"age":             &types.AttributeValueMemberS{Value: task.Age},
		"region":          &types.AttributeValueMemberS{Value: task.Region},
		"author_id":       &types.AttributeValueMemberS{Value: task.AuthorId},
		"article_date_id": &types.AttributeValueMemberS{Value: task.ArticleDateId},
		"valid":           &types.AttributeValueMemberBOOL{Value: task.Valid},
	}, nil
}

func (task *Task) Serialize(av map[string]types.AttributeValue) {
	task.ID = av["id"].(*types.AttributeValueMemberS).Value
	task.Prompt = av["prompt"].(*types.AttributeValueMemberS).Value
	task.Age = av["age"].(*types.AttributeValueMemberS).Value
	task.Region = av["region"].(*types.AttributeValueMemberS).Value
	task.AuthorId = av["author_id"].(*types.AttributeValueMemberS).Value
	task.ArticleDateId = av["article_date_id"].(*types.AttributeValueMemberS).Value
	task.Valid = av["valid"].(*types.AttributeValueMemberBOOL).Value
}
