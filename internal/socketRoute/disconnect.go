package socketRoute

import (
	"context"
	"hyper_api/internal/dto"
	"hyper_api/internal/utils/aws"
)

func Disconnect(ctx context.Context, event dto.EventBody) error {
	svc, err := aws.NewDynamoDBClient()
	if err != nil {
		return err
	}
	connectionId := ctx.Value("connectionId").(string)
	err = aws.UnsubscribeFromDB(svc, connectionId)
	if err != nil {
		return err
	}
	return nil
}
