package socketRoute

import (
	"context"
	"encoding/json"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	apigw "github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/aws"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	awsUtil "hyper_api/internal/utils/aws"
)

type Task struct {
	TaskId string `json:"taskID"`
	Status string `json:"status"`
}

func TaskDone(ctx context.Context, event dto.EventBody) error {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	c := config.GetConfig()
	if err != nil {
		return err
	}
	client := apigw.NewFromConfig(cfg, func(o *apigw.Options) {
		o.EndpointResolver = apigw.EndpointResolverFromURL(c.WebSocketURL)
	})
	svc, err := awsUtil.NewDynamoDBClient()
	if err != nil {
		return err
	}
	var task Task
	err = json.Unmarshal([]byte(event.Body), &task)
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
	}
	return nil
}
