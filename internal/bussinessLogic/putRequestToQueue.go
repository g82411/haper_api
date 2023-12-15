package bussinessLogic

import (
	"fmt"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws"
)

func PutImageRequestToQueue(taskId, userId, prompt, articleId string) error {
	setting := config.GetConfig()
	task := dto.GenerateImageTask{
		TaskID:    taskId,
		AuthorId:  userId,
		Prompt:    prompt,
		ArticleId: articleId,
	}
	sqsClient, err := aws.NewSQSClient(setting.SQSQueueUrl)
	if err != nil {
		return fmt.Errorf("error creating sqs client %v", err)
	}
	err = sqsClient.SendJSONMessage(task)
	if err != nil {
		return fmt.Errorf("error sending message %v", err)
	}
	return nil
}
