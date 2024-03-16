package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"hyper_api/internal/dto"
)

type CognitoClient struct {
	cip        *cognitoidentityprovider.CognitoIdentityProvider
	userPoolId string
}

func NewCognitoClient() (*CognitoClient, error) {
	sess, err := NewAWSSession()
	if err != nil {
		return nil, err
	}

	cip := cognitoidentityprovider.New(sess)
	return &CognitoClient{cip: cip}, nil
}

func NewAdminCognitoClient(userPoolId string) (*CognitoClient, error) {
	client, err := NewCognitoClient()
	if err != nil {
		return nil, err
	}
	client.userPoolId = userPoolId
	return client, nil
}

func (client *CognitoClient) GetUserAttributeByAccessToken(accessToken string) (*dto.UserInfo, error) {
	var userInfo dto.UserInfo
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: &accessToken,
	}

	result, err := client.cip.GetUser(input)
	if err != nil {
		return &userInfo, err
	}
	userInfo.InternalUserName = *result.Username
	for _, v := range result.UserAttributes {
		if *v.Name == "sub" {
			userInfo.Sub = *v.Value
		}
		if *v.Name == "name" {
			userInfo.Name = *v.Value
		}
		if *v.Name == "email" {
			userInfo.Email = *v.Value
		}
		if *v.Name == "custom:isDoneSurvey" {
			userInfo.IsDoneSurvey = *v.Value
		}
		if *v.Name == "picture" {
			userInfo.Picture = *v.Value
		}
	}

	return &userInfo, nil
}

func (client *CognitoClient) UpdateUserInfo(usernameInCognito, name, age, career, gender string) error {
	userAttributes := []*cognitoidentityprovider.AttributeType{
		{
			Name:  aws.String("custom:real_name"),
			Value: aws.String(name),
		},
		{
			Name:  aws.String("custom:age"),
			Value: aws.String(age),
		},
		{
			Name:  aws.String("custom:career"),
			Value: aws.String(career),
		},
		{
			Name:  aws.String("gender"),
			Value: aws.String(gender),
		},
		{
			Name:  aws.String("custom:isDoneSurvey"),
			Value: aws.String("true"),
		},
	}

	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     &client.userPoolId,
		Username:       &usernameInCognito,
		UserAttributes: userAttributes,
	}

	_, err := client.cip.AdminUpdateUserAttributes(input)
	return err
}
