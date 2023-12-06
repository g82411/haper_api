package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"hyper_api/internal/config"
)

func UpdateUserSurveyStatus(userSub string, isDoneSurvey bool) error {
	sess, err := NewAWSSession()
	if err != nil {
		return err
	}

	cip := cognitoidentityprovider.New(sess)
	c := config.GetConfig()
	userPoolId := c.CognitoUserPool

	value := "false"
	if isDoneSurvey {
		value = "true"
	}

	userAttributes := []*cognitoidentityprovider.AttributeType{
		{
			Name:  aws.String("custom:isDoneSurvey"),
			Value: aws.String(value),
		},
	}

	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     &userPoolId,
		Username:       &userSub,
		UserAttributes: userAttributes,
	}

	_, err = cip.AdminUpdateUserAttributes(input)
	return err
}
