package bussinessLogic

import (
	"context"
	"fmt"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws/dynamodb"
)

const PromptTemplate = "你是一位懂得easyread圖片規則的設計師，了解心智障礙者閱讀資訊需要簡單易懂的圖面資訊，請 %s \n\n產出的圖片將符合以下幾個重點: \n1.  圖片以卡通插圖的風格生成，線條乾淨俐落，線條較粗，避免複雜、多餘的線條，使用簡單的色彩、元素設計\n2. 產出的圖片內容會是簡單，容易看懂\n3. 圖片僅有主體、呈現完整的樣貌、貼近實際場景、貼近實際動作。\n4.主體為彩色，圖片背景是白色，圖片內容限縮在1024x1024以內，保留主體完整性不要切邊。\n\n"

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
