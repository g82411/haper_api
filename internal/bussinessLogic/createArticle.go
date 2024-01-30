package bussinessLogic

import (
	"context"
	"fmt"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws/dynamodb"
	"hyper_api/internal/utils/resolver/date"
)

func CreateArticle(ctx context.Context, authorId, authorName, keyword string) (*models.Article, error) {
	id := utils.GenerateShortKey()
	today := date.GetNowDateString()
	article := models.Article{
		AuthorId:   authorId,
		AuthorName: authorName,
		Keyword:    keyword,
		Valid:      false,
		ID:         id,
		Date:       today,
	}

	err := dynamodb.Insert(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("error while inserting article: %s", err.Error())
	}
	return &article, nil
}
