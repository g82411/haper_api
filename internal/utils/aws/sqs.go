package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"hyper_api/internal/config"
)

type SQSClient struct {
	session  *session.Session
	client   *sqs.SQS
	queueUrl string
}

func NewSQSClient() (*SQSClient, error) {
	setting := config.GetConfig()
	region := setting.AWSRegion
	queueUrl := setting.SQSQueueURL

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), // 替换为你的 AWS 区域
	})
	if err != nil {
		return nil, err
	}
	return &SQSClient{
		session:  sess,
		client:   sqs.New(sess),
		queueUrl: queueUrl,
	}, nil
}

func encodeMessageToJSON(message interface{}) (string, error) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(msgBytes), nil
}

func (s *SQSClient) SendJSONMessage(message interface{}) error {
	jsonMessage, err := encodeMessageToJSON(message)
	if err != nil {
		return err
	}
	_, err = s.client.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(jsonMessage),
		QueueUrl:    aws.String(s.queueUrl),
	})
	if err != nil {
		return err
	}
	return nil
}
