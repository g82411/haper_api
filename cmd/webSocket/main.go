package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hyper_api/internal/dto"
	lambdaDto "hyper_api/internal/dto/lambda"
	"hyper_api/internal/socketRoute"
)

func buildContext(ctx context.Context, event lambdaDto.WebSocketEvent) context.Context {
	nextCtx := context.WithValue(ctx, "connectionId", event.RequestContext.ConnectionID)
	return nextCtx
}

func Handler(ctx context.Context, event lambdaDto.WebSocketEvent) error {
	body := event.Body
	routeKey := event.RequestContext.RouteKey
	var socketEvent dto.EventBody
	err := json.Unmarshal([]byte(body), &event)
	nextCtx := buildContext(ctx, event)
	if err != nil {
		return fmt.Errorf("error parsing message body %v", err)
	}
	if routeKey == "$disconnect" {
		err = socketRoute.Disconnect(nextCtx, socketEvent)
	}
	if socketEvent.Action == "subscribe" {
		err = socketRoute.Subscribe(nextCtx, socketEvent)
	}
	if socketEvent.Action == "taskDone" {
		err = socketRoute.TaskDone(nextCtx, socketEvent)
	}
	return err
}

func main() {
	lambda.Start(Handler)
}
