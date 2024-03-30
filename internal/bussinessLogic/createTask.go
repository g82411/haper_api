package bussinessLogic

import (
	"context"
	"fmt"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws/dynamodb"
)

const PromptTemplate = `
請繪製：
閱讀對象：輕度智能障礙者
閱讀對象年齡：%s
閱讀對象地區：%s
預計搭配的說明文字：%s

以下是風格條件說明：
1.完整性：在圖片大小下，完整呈現圖片，不會畫出不完整的畫面
2.簡潔性：
   2-1. 圖片底色為白色，不會有裝飾邊、色塊等
   2-2. 不使用複雜的配色
   2-3. 除了主體，不使用其他元素
3.具象性：不使用特殊的標誌或抽象的圖片
指導原則： 你會遵照以下原則繪製圖片：
1. 風格：卡通插圖
2. 線條：寬度為2 pixel以上，且單一線條，不花俏
3. 色彩：選擇高對比度的色彩組合來提高可讀性，並使用色彩來強化或區分信息。避免使用過於鮮豔或刺眼的顏色，這可能會分散注意力或對某些閱讀對象造成不適。
4. 元素：繪製清晰、易於識別的圖示和符號。圖片應該一目了然，避免多重含義或容易引起誤解的元素。
5. 內容：圖片內容應與預計搭配的說明文字直接相關，幫助讀者理解文本的具體內容。
6. 內容文化：符合閱讀對象地區文化與實際情況。
7. 閱讀對象：產生適合閱讀對象的圖片，例如當閱讀對象為大人時，應避免使用適合兒童的圖片。

限制：
1.完整性：在圖片大小下，完整呈現圖片，不會畫出不完整的畫面
2.簡潔性：
   2-1. 圖片底色為白色，不會有裝飾邊、色塊等
   2-2. 不使用複雜的配色
   2-3. 除了主體，不使用其他元素
3.具象性：不使用特殊的標誌或抽象的圖片
`

func CreateTask(ctx context.Context, article *models.Article) (*models.Task, error) {
	prompt := fmt.Sprintf(PromptTemplate, article.Age, article.Region, article.Keyword)
	task := models.Task{
		ID:            utils.GenerateShortKey(),
		ArticleDateId: article.DateId,
		AuthorId:      article.AuthorId,
		Prompt:        prompt,
		Age:           article.Age,
		Region:        article.Region,
		Valid:         true,
	}
	err := dynamodb.Insert(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("error while inserting article: %s", err.Error())
	}
	return &task, nil
}
