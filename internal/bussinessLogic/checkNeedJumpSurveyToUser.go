package bussinessLogic

import "hyper_api/internal/models"

func CheckNeedJumpSurveyToUser(sub string) (bool, error) {
	ArticleCountThatNeedSurvey := [3]int{5, 20, 40}

	dbClient, err := models.NewDBClient()
	if err != nil {
		return false, err
	}
	var count int64
	tx := dbClient.Model(&models.Article{})
	tx.Where("author_id = ?", sub)
	tx.Where("valid = ?", true)
	tx.Limit(ArticleCountThatNeedSurvey[2] + 1)
	tx.Count(&count)

	flag := false
	for i := 0; i < 3; i++ {
		if count == int64(ArticleCountThatNeedSurvey[i]) {
			flag = true
			break
		}
	}
	return flag, nil
}
