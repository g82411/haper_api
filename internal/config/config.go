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
	DBHost     string
	DBPassword string
	DBUsername string
	DBName     string
}

func GetConfig() *Config {
	return &Config{
		AWSRegion:       "ap-southeast-1",
		SQSQueueName:    os.Getenv("SQS_QUEUE_NAME"),
		SQSQueueURL:     os.Getenv("SQS_QUEUE_URL"),
		OpenAIKey:       os.Getenv("OPENAI_KEY"),
		CognitoDomain:   os.Getenv("COGNITO_DOMAIN"),
		CognitoClientId: os.Getenv("COGNITO_CLIENT_ID"),
		RedirectURL:     os.Getenv("REDIRECT_URL"),
		TokenTable:      os.Getenv("TOKEN_TABLE"),
		SurveyTable:     os.Getenv("SURVEY_TABLE"),
		CognitoUserPool: os.Getenv("COGNITO_USER_POOL"),
		AllowOrigin:     os.Getenv("ALLOW_ORIGIN"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBUsername:      os.Getenv("DB_USERNAME"),
		DBName:          os.Getenv("DB_NAME"),
	}
}
