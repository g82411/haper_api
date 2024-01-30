package models

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Task struct {
	ID        string
	Prompt    string
	AuthorId  string
	ArticleId string
	Valid     bool
}

func (task *Task) TableName(ctx context.Context) string {
	stage := ctx.Value("stage").(string)
	return fmt.Sprintf("%s_tasks", stage)
}

func (task *Task) Deserialize() (map[string]types.AttributeValue, error) {
	return map[string]types.AttributeValue{
		"id":         &types.AttributeValueMemberS{Value: task.ID},
		"prompt":     &types.AttributeValueMemberS{Value: task.Prompt},
		"author_id":  &types.AttributeValueMemberS{Value: task.AuthorId},
		"article_id": &types.AttributeValueMemberS{Value: task.ArticleId},
		"valid":      &types.AttributeValueMemberBOOL{Value: task.Valid},
	}, nil
}

func (task *Task) Serialize(av map[string]types.AttributeValue) {
	task.ID = av["id"].(*types.AttributeValueMemberS).Value
	task.Prompt = av["prompt"].(*types.AttributeValueMemberS).Value
	task.AuthorId = av["author_id"].(*types.AttributeValueMemberS).Value
	task.ArticleId = av["article_id"].(*types.AttributeValueMemberS).Value
	task.Valid = av["valid"].(*types.AttributeValueMemberBOOL).Value
}
