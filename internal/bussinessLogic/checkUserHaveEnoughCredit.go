package bussinessLogic

import (
	"fmt"
	"hyper_api/internal/models"
	"time"
)

const DailyCredit = 10

func CheckUserHaveEnoughCredit(sub string) (bool, error) {
	dbClient, err := models.NewDBClient()
	if err != nil {
		return false, fmt.Errorf("error when connect to db %v", err)
	}
	var count int64
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Unix()
	tx := dbClient.Model(&models.Article{})
	tx.Where("author_id = ?", sub)
	tx.Where("valid = ?", true)
	tx.Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay)
	tx.Limit(DailyCredit + 1)
	tx.Count(&count)
	if count > DailyCredit {
		return false, nil
	}
	return true, nil
}
