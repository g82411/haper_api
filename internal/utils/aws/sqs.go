package aws

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	client   *sqs.Client
	queueUrl string
}

func NewSQSClient(queueUrl string) (*SQSClient, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// 创建 SQS 客户端
	client := sqs.NewFromConfig(cfg)
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// 创建 SQS 客户端
	if err != nil {
		return nil, err
	}
	return &SQSClient{
		client:   client,
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
	_, err = s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(jsonMessage),
		QueueUrl:    aws.String(s.queueUrl),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *SQSClient) PullMessages() (*sqs.ReceiveMessageOutput, error) {
	result, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(s.queueUrl),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SQSClient) DeleteMessage(message types.Message) error {
	_, err := s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return err
	}
	return nil
}
