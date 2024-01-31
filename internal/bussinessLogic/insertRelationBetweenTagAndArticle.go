package bussinessLogic

import (
	"context"
	"fmt"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws/dynamodb"
)

func InsertRelationBetweenTagAndArticle(ctx context.Context, article *models.Article, tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	index := make([]dynamodb.SerializeAble, 0)
	invertedIndex := make([]dynamodb.SerializeAble, 0)
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		index = append(index, models.ArticleTag{
			ArticleID: article.ID,
			TagName:   tag,
		})
		invertedIndex = append(invertedIndex, models.TagArticle{
			ArticleID: article.ID,
			TagName:   tag,
		})
	}
	if len(index) == 0 {
		return nil
	}
	err := dynamodb.BulkInsert(ctx, index)
	if err != nil {
		return fmt.Errorf("insert relation between tag and article: %v", err)
	}
	err = dynamodb.BulkInsert(ctx, invertedIndex)
	if err != nil {
		return fmt.Errorf("insert relation between tag and article: %v", err)
	}
	return nil
}
