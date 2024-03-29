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
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"hyper_api/internal/config"
	"hyper_api/internal/middleware"
	"hyper_api/internal/routes"
	"io"
	"net/http"
	"net/url"
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
	urlStr := req.Path // 假设 req.Path 包含基本的 URL 路径
	if len(req.QueryStringParameters) > 0 {
		queryParams := url.Values{}
		for key, value := range req.QueryStringParameters {
			queryParams.Add(key, value)
		}
		// 将查询参数附加到 URL
		urlStr = fmt.Sprintf("%s?%s", urlStr, queryParams.Encode())
	}
	// 创建新的 HTTP 请求
	httpReq, err := http.NewRequest(req.HTTPMethod, urlStr, bytes.NewBuffer(jsonStr))
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
		"Access-Control-Allow-Credentials": "true",
	}
	// 初始化 Fiber
	app := fiber.New()
	log.SetLevel(log.LevelInfo)

	// 定义路由
	setUpApp(app)
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
	for key, value := range httpRes.Header {
		CommonHeader[key] = value[0] // 假设每个header只有一个值
	}
	fmt.Println("Response Header: ", CommonHeader)
	return events.APIGatewayProxyResponse{
		StatusCode: httpRes.StatusCode,
		Headers:    CommonHeader,
		Body:       string(r),
	}, nil
}

func setUpApp(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("env:isLocal", lambdacontext.FunctionName == "")
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowHeaders:     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(middleware.LoadConfig)
	routes.BindingRoutes(app)

}

func main() {
	// 判断是否在 Lambda 环境中运行
	if lambdacontext.FunctionName != "" {
		lambda.Start(handler)
	} else {
		// 作为独立的 HTTP 服务运行
		app := fiber.New()
		log.SetLevel(log.LevelDebug)
		app.Use(cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool { return true },
			AllowHeaders:     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		}))
		setUpApp(app)
		log.Fatal(app.Listen(":8080"))
	}
}
