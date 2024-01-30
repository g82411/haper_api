package bussinessLogic

import (
	"fmt"
	"hyper_api/internal/config"
	"hyper_api/internal/models"
	"hyper_api/internal/utils/aws"
)

func PutImageRequestToQueue(task *models.Task) error {
	setting := config.GetConfig()
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
