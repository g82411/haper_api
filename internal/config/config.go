package config

import "os"

type Config struct {
	AWSRegion string
	// SQS
	RedirectURL     string
	CognitoDomain   string
	SQSQueueName    string
	SQSQueueURL     string
	OpenAIKey       string
	TokenTable      string
	SurveyTable     string
	SubscriberTable string
	CognitoUserPool string
	AllowOrigin     string
	// DB
	DBHost           string
	DBPassword       string
	DBUsername       string
	DBName           string
	CDNHost          string
	GenerateS3Bucket string
	WebHOST          string
	// web socket
	WebSocketHost string
	WebSocketPath string
	WebSocketURL  string
	SQSQueueUrl   string
	Env           string
}

func GetConfig() *Config {
	return &Config{
		AWSRegion:        "ap-southeast-1",
		WebHOST:          os.Getenv("WEB_HOST"),
		OpenAIKey:        os.Getenv("OPENAI_KEY"),
		SQSQueueName:     os.Getenv("SQS_QUEUE_NAME"),
		SQSQueueURL:      os.Getenv("SQS_QUEUE_URL"),
		CognitoUserPool:  os.Getenv("COGNITO_USER_POOL"),
		TokenTable:       os.Getenv("TOKEN_TABLE"),
		SurveyTable:      os.Getenv("SURVEY_TABLE"),
		SubscriberTable:  os.Getenv("SUBSCRIBER_TABLE"),
		AllowOrigin:      os.Getenv("ALLOW_ORIGIN"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBUsername:       os.Getenv("DB_USERNAME"),
		DBName:           os.Getenv("DB_NAME"),
		RedirectURL:      os.Getenv("REDIRECT_URL"),
		CDNHost:          os.Getenv("CDN_HOST"),
		GenerateS3Bucket: os.Getenv("GENERATE_S3_BUCKET"),
		WebSocketHost:    os.Getenv("WEB_SOCKET_HOST"),
		WebSocketPath:    os.Getenv("WEB_SOCKET_PATH"),
		WebSocketURL:     os.Getenv("WEB_SOCKET_URL"),
		SQSQueueUrl:      os.Getenv("SQS_QUEUE_URL"),
		Env:              os.Getenv("ENV"),
	}
}
