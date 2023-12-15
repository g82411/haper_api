package socketRoute

import (
	"context"
	"encoding/json"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws"
)

type SubscribeTask struct {
	TaskId string `json:"taskID"`
}

func Subscribe(ctx context.Context, event dto.EventBody) error {
	svc, err := aws.NewDynamoDBClient()
	if err != nil {
		return err
	}
	var subscribeTask SubscribeTask
	connectionId := ctx.Value("connectionId").(string)
	err = json.Unmarshal([]byte(event.Body), &subscribeTask)
	if err != nil {
		return err
	}
	err = aws.PutSubscriberToDB(svc, connectionId, subscribeTask.TaskId)
	if err != nil {
		return err
	}
	return nil
}
