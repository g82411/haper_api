package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"hyper_api/internal/config"
	"hyper_api/internal/routes"
	"io"
	"log"
	"net/http"
)

func APIGatewayRequestToHTTPRequest(req events.APIGatewayProxyRequest) (*http.Request, error) {
	// 解析 JSON body (如果需要)
	var jsonStr = []byte(req.Body)
	if req.IsBase64Encoded {
		// 如果 body 是 base64 编码的，先进行解码
		decodedBody, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return nil, err
		}
		jsonStr = decodedBody
	}

	// 创建新的 HTTP 请求
	httpReq, err := http.NewRequest(req.HTTPMethod, req.Path, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	// 复制请求头
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}

// 定义 Lambda Handler

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := config.GetConfig()
	CommonHeader := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
		"Access-Control-Allow-Methods":     "OPTIONS,GET,PUT,POST,DELETE,PATCH,HEAD",
		"Access-Control-Allow-Origin":      c.AllowOrigin,
		"Access-Control-Allow-Credentials": "false",
	}
	// 初始化 Fiber
	app := fiber.New()

	// 定义路由
	setupRoutes(app)
	fmt.Printf("Request Body: %v\n", req.Body)
	fmt.Printf("Request Header: %v\n", req.Headers)

	httpReq, err := APIGatewayRequestToHTTPRequest(req)
	if err != nil {
		fmt.Printf("Error when parsing request %v", err)
		return events.APIGatewayProxyResponse{Headers: CommonHeader}, err
	}

	// 使用 Fiber 处理请求
	httpRes, err := app.Test(httpReq, 12000*1000)
	if err != nil {
		fmt.Printf("Error when deal request %v", err)
		return events.APIGatewayProxyResponse{Headers: CommonHeader}, err
	}

	// 将 HTTP 响应转换为 API Gateway 响应
	if err != nil {
		fmt.Printf("Error when parsing response %v", err)
		return events.APIGatewayProxyResponse{Headers: CommonHeader}, err
	}
	r, _ := io.ReadAll(httpRes.Body)
	return events.APIGatewayProxyResponse{
		StatusCode: httpRes.StatusCode,
		Headers:    CommonHeader,
		Body:       string(r),
	}, nil
}

func setupRoutes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Message": "pong"})
	})
	routes.BindingRoutes(app)
}

func main() {
	// 判断是否在 Lambda 环境中运行
	if lambdacontext.FunctionName != "" {
		lambda.Start(handler)
	} else {
		// 作为独立的 HTTP 服务运行
		// TODO: refresh middle ware
		app := fiber.New()
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:3000",
			AllowHeaders:     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		}))
		setupRoutes(app)
		log.Fatal(app.Listen(":8080"))
	}
}
