package socketRoute

import (
	"context"
	"encoding/json"
	"fmt"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	apigw "github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/aws"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	awsUtil "hyper_api/internal/utils/aws"
)

type Task struct {
	TaskId    string `json:"taskID"`
	ArticleId string `json:"articleID"`
	Status    string `json:"status"`
}

func TaskDone(ctx context.Context, event dto.EventBody) error {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	c := config.GetConfig()
	if err != nil {
		return err
	}
	client := apigw.NewFromConfig(cfg, func(o *apigw.Options) {
		o.BaseEndpoint = aws.String(c.WebSocketURL)
	})
	svc, err := awsUtil.NewDynamoDBClient()
	if err != nil {
		return err
	}
	var task Task
	err = json.Unmarshal([]byte(event.Body), &task)
	fmt.Printf("task: %v\n", task)
	if err != nil {
		return err
	}
	subscribers, err := awsUtil.GetSubscriberFromDB(svc, task.TaskId)
	if err != nil {
		return err
	}
	for _, subscriber := range subscribers {
		// 发送消息
		_, err = client.PostToConnection(context.TODO(), &apigw.PostToConnectionInput{
			ConnectionId: aws.String(subscriber),
			Data:         []byte(event.Body),
		})
		if err != nil {
			return fmt.Errorf("error sending message %v", err)
		}
	}
	return nil
}
