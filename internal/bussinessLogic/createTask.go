package bussinessLogic

import (
	"context"
	"fmt"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws/dynamodb"
)

const PromptTemplate = "%v以卡通插圖的風格繪製，線條乾淨俐落，線條較粗，避免複雜、多餘的線條，使用簡單的色彩。\n圖片僅有主體、呈現完整的樣貌、貼近實際場景、貼近實際動作。主體為彩色，圖片背景是白色。"

func CreateTask(ctx context.Context, article *models.Article) (*models.Task, error) {
	prompt := fmt.Sprintf(PromptTemplate, article.Keyword)
	task := models.Task{
		ID:            utils.GenerateShortKey(),
		ArticleDateId: article.DateId,
		AuthorId:      article.AuthorId,
		Prompt:        prompt,
		Valid:         true,
	}
	err := dynamodb.Insert(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("error while inserting article: %s", err.Error())
	}
	return &task, nil
}
