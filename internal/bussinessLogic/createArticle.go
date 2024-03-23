package bussinessLogic

import (
	"context"
	"fmt"
	"hyper_api/internal/dto"
	"hyper_api/internal/models"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws/dynamodb"
	"hyper_api/internal/utils/resolver/date"
)

func CreateArticle(ctx context.Context, userInfo *dto.UserInfo, keyword string, age string, region string) (*models.Article, error) {
	id := utils.GenerateShortKey()
	today := date.GetNowDateString()
	article := models.Article{
		AuthorId:    userInfo.Sub,
		AuthorName:  userInfo.Name,
		Keyword:     keyword,
		Valid:       false,
		ID:          id,
		Date:        today,
		DateId:      today + "_" + id,
		AuthorImage: userInfo.Picture,
		Age:         age,
		Region:      region,
	}
	err := dynamodb.Insert(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("error while inserting article: %s", err.Error())
	}
	return &article, nil
}
