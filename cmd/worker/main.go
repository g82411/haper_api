package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"hyper_api/internal/bussinessLogic"
	"hyper_api/internal/dto"
	"hyper_api/internal/notify"
	"hyper_api/internal/tasks"
)

func parseMessage(message string) (dto.GenerateImageTask, error) {
	var task dto.GenerateImageTask
	err := json.Unmarshal([]byte(message), &task)
	if err != nil {
		return dto.GenerateImageTask{}, fmt.Errorf("error parsing message body %v", err)
	}
	return task, nil
}

func runTask(task dto.GenerateImageTask) error {
	err := tasks.GenerateImageToOpenAI(task)
	if err != nil {
		return fmt.Errorf("error running task %v", err)
	}
	return nil
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		task, err := parseMessage(message.Body)
		err = runTask(task)
		if err != nil {
			return fmt.Errorf("error running task %v", err)
		}
		err = notify.CompleteTask(task)
		if err != nil {
			return fmt.Errorf("error completing task %v", err)
		}
	}
	return nil
}

func localHandler() error {
	queueTasks, err := bussinessLogic.PullRequestsFromQueue()
	if err != nil {
		return fmt.Errorf("error parsing message body %v", err)
	}
	for _, task := range queueTasks {
		err := runTask(task)
		if err != nil {
			return fmt.Errorf("error running task %v", err)
		}
		err = notify.CompleteTask(task)
		if err != nil {
			return fmt.Errorf("error completing task %v", err)
		}
	}
	return nil
}

func main() {
	if lambdacontext.FunctionName != "" {
		lambda.Start(handler)
		return
	}
	err := localHandler()
	if err != nil {
		fmt.Printf("error running local handler %v", err)
	}
}
