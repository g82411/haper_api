package config

import "os"

type Config struct {
	AWSRegion string
	// SQS
	CognitoClientId string
	RedirectURL     string
	CognitoDomain   string
	SQSQueueName    string
	SQSQueueURL     string
	OpenAIKey       string
	TokenTable      string
	SurveyTable     string
	CognitoUserPool string
	AllowOrigin     string
	// DB
	DBHost           string
	DBPassword       string
	DBUsername       string
	DBName           string
	CDNHost          string
	GenerateS3Bucket string
}

func GetConfig() *Config {
	return &Config{
		AWSRegion:        "ap-southeast-1",
		OpenAIKey:        os.Getenv("OPENAI_KEY"),
		SQSQueueName:     os.Getenv("SQS_QUEUE_NAME"),
		SQSQueueURL:      os.Getenv("SQS_QUEUE_URL"),
		CognitoDomain:    os.Getenv("COGNITO_DOMAIN"),
		CognitoClientId:  os.Getenv("COGNITO_CLIENT_ID"),
		CognitoUserPool:  os.Getenv("COGNITO_USER_POOL"),
		TokenTable:       os.Getenv("TOKEN_TABLE"),
		SurveyTable:      os.Getenv("SURVEY_TABLE"),
		AllowOrigin:      os.Getenv("ALLOW_ORIGIN"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBUsername:       os.Getenv("DB_USERNAME"),
		DBName:           os.Getenv("DB_NAME"),
		RedirectURL:      os.Getenv("REDIRECT_URL"),
		CDNHost:          os.Getenv("CDN_HOST"),
		GenerateS3Bucket: os.Getenv("GENERATE_S3_BUCKET"),
	}
}
