package aws

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"hyper_api/internal/config"
	"hyper_api/internal/utils"
)

func NewDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func GetTokenFromDB(svc *dynamodb.Client, userID string) (string, string, error) {
	c := config.GetConfig()
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.TokenTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := svc.GetItem(context.Background(), input)
	if err != nil {
		return "", "", err
	}

	return result.Item["IDToken"].(*types.AttributeValueMemberS).Value, result.Item["RefreshToken"].(*types.AttributeValueMemberS).Value, nil
}

func PutTokenToDB(svc *dynamodb.Client, userID, idToken, refreshToken string) error {
	c := config.GetConfig()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.TokenTable),
		Item: map[string]types.AttributeValue{
			"id":           &types.AttributeValueMemberS{Value: userID},
			"IDToken":      &types.AttributeValueMemberS{Value: idToken},
			"RefreshToken": &types.AttributeValueMemberS{Value: refreshToken},
		},
	}

	_, err := svc.PutItem(context.TODO(), input)
	return err
}

func PutSurveyResultToDB(svc *dynamodb.Client, userID, surveyResult string) error {
	c := config.GetConfig()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.SurveyTable),
		Item: map[string]types.AttributeValue{
			"id":           &types.AttributeValueMemberS{Value: userID},
			"surveyResult": &types.AttributeValueMemberS{Value: surveyResult},
		},
	}
	_, err := svc.PutItem(context.TODO(), input)
	return err
}

func GetUserInfoByAccessToken(svc *dynamodb.Client, accessToken string) (utils.Claims, error) {
	var ret utils.Claims
	accessData, err := utils.ExtractUserInfoFromToken(accessToken)
	if err != nil {
		return ret, err
	}
	userSub := accessData.Sub
	idToken, _, err := GetTokenFromDB(svc, userSub)
	if err != nil {
		return ret, err
	}
	userData, err := utils.ExtractUserInfoFromToken(idToken)
	if err != nil {
		return ret, err
	}
	return userData, nil
}

func PutSubscriberToDB(svc *dynamodb.Client, subscriber, taskId string) error {
	c := config.GetConfig()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.SubscriberTable),
		Item: map[string]types.AttributeValue{
			"id":     &types.AttributeValueMemberS{Value: subscriber},
			"taskId": &types.AttributeValueMemberS{Value: taskId},
		},
	}
	_, err := svc.PutItem(context.Background(), input)
	return err
}

func GetSubscriberFromDB(svc *dynamodb.Client, taskId string) ([]string, error) {
	c := config.GetConfig()
	input := &dynamodb.ScanInput{
		TableName: aws.String(c.SubscriberTable),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":taskId": &types.AttributeValueMemberS{Value: taskId},
		},
		FilterExpression: aws.String("taskId = :taskId"),
	}
	result, err := svc.Scan(context.Background(), input)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, item := range result.Items {
		ret = append(ret, item["id"].(*types.AttributeValueMemberS).Value)
	}
	return ret, nil
}

func UnsubscribeFromDB(svc *dynamodb.Client, subscriber string) error {
	c := config.GetConfig()
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.SubscriberTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: subscriber},
		},
	}
	_, err := svc.DeleteItem(context.Background(), input)
	return err
}
