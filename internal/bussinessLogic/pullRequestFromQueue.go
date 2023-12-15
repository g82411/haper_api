package bussinessLogic

import (
	"encoding/json"
	"fmt"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws"
)

func parseMessage(message string) (dto.GenerateImageTask, error) {
	var task dto.GenerateImageTask
	err := json.Unmarshal([]byte(message), &task)
	if err != nil {
		return dto.GenerateImageTask{}, fmt.Errorf("error parsing message body %v", err)
	}
	return task, nil
}

func PullRequestsFromQueue() ([]dto.GenerateImageTask, error) {
	setting := config.GetConfig()
	var tasks []dto.GenerateImageTask
	sqsClient, err := aws.NewSQSClient(setting.SQSQueueUrl)
	if err != nil {
		return tasks, err
	}
	messages, err := sqsClient.PullMessages()
	if err != nil {
		return tasks, err
	}
	for _, message := range messages.Messages {
		task, err := parseMessage(*message.Body)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
		go sqsClient.DeleteMessage(message)
	}
	return tasks, nil
}
