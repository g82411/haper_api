package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewAWSSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"), // 替换为你的 AWS 区域
	})
}
