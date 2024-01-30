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
	n := len(tags)
	index := make([]dynamodb.SerializeAble, n)
	invertedIndex := make([]dynamodb.SerializeAble, n)
	for i, tag := range tags {
		index[i] = models.ArticleTag{
			ArticleID: article.ID,
			TagName:   tag,
		}
		invertedIndex[i] = models.TagArticle{
			ArticleID: article.ID,
			TagName:   tag,
		}
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
