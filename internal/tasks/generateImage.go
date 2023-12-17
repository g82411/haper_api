package tasks

import (
	"fmt"
	"gorm.io/gorm"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/models"
)

func markTaskAsDone(db *gorm.DB, taskId string) error {
	var task models.Task
	tx := db.Debug().Table("tasks").Where("id = ?", taskId).Find(&task)
	if tx.Error != nil {
		return tx.Error
	}
	if task.ID == "" {
		return nil
	}
	task.Status = 1
	tx = db.Debug().Table("tasks").Where("id = ?", taskId).Save(&task)
	return nil
}

func saveImageToArticle(db *gorm.DB, articleId, imageUrl string) error {
	var article models.Article
	tx := db.Debug().Table("articles").Where("id = ?", articleId).Find(&article)
	if tx.Error != nil {
		return tx.Error
	}
	if article.ID == "" {
		return nil
	}
	article.Valid = true
	article.Url = imageUrl
	tx = db.Debug().Table("articles").Where("id = ?", articleId).Save(&article)
	return nil
}

func GenerateImageToOpenAI(task dto.GenerateImageTask) error {
	generatedImage, err := bussinessLogic.GenerateImageByPrompt(task.Prompt)
	if err != nil {
		return fmt.Errorf("generate Image error %v", err)
	}

	dbClient, err := models.NewDBClient()
	if err != nil {
		return fmt.Errorf("error when establish db connection %v", err)
	}
	err = markTaskAsDone(dbClient, task.TaskID)
	if err != nil {
		return fmt.Errorf("error when mark as complete %v", err)
	}
	err = saveImageToArticle(dbClient, task.ArticleId, generatedImage)
	if err != nil {
		return fmt.Errorf("error when save image to article %v", err)
	}
	return nil
}
