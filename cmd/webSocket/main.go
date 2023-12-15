package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hyper_api/internal/dto"
	lambdaDto "hyper_api/internal/dto/lambda"
	"hyper_api/internal/socketRoute"
)

func buildContext(ctx context.Context, event lambdaDto.WebSocketEvent) context.Context {
	nextCtx := context.WithValue(ctx, "connectionId", event.RequestContext.ConnectionID)
	return nextCtx
}

func handleWebSocketRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Handling WebSocket request for route: %s\n", event.RequestContext.RouteKey)

	// 根据不同的 routeKey 执行不同的逻辑
	switch event.RequestContext.RouteKey {
	case "$connect":
		// 处理连接事件
		fmt.Println("Connected:", event.RequestContext.ConnectionID)
	case "$disconnect":
		// 处理断开连接事件
		fmt.Println("Disconnected:", event.RequestContext.ConnectionID)
	case "$default":
		// 处理默认路由（接收消息）
		var message map[string]interface{}
		err := json.Unmarshal([]byte(event.Body), &message)
		if err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			return events.APIGatewayProxyResponse{}, err
		}
		fmt.Printf("Received message: %v\n", message)
	}

	// 返回响应
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
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
	lambda.Start(handleWebSocketRequest)
}
