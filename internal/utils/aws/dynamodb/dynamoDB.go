package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type InputQuery struct {
	ExpressionAttribute    *map[string]types.AttributeValue
	FilterExpression       *string
	KeyConditionExpression *string
	ExclusiveStartKey      *map[string]types.AttributeValue
	IndexName              string
	ScanIndexForward       bool
	Limit                  int32
}

type UpdateQuery struct {
	ExpressionAttribute *map[string]types.AttributeValue
	UpdateExpression    *string
	Key                 *map[string]types.AttributeValue
	Limit               int32
	ReturnValues        types.ReturnValue
}

type SerializeAble interface {
	Serialize(av map[string]types.AttributeValue) interface{}
	Deserialize() map[string]types.AttributeValue
	TableName(ctx context.Context) string
}

func WithDynamoDBConnection(ctx context.Context) (context.Context, error) {
	if ctx.Value("dynamodb") != nil {
		return ctx, nil
	}
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	tx := dynamodb.NewFromConfig(cfg)
	ctx = context.WithValue(ctx, "dynamodb", tx)
	return ctx, nil
}

func Insert(ctx context.Context, data SerializeAble) error {
	svc := ctx.Value("dynamodb").(*dynamodb.Client)
	av := data.Deserialize()
	tableName := data.TableName(ctx)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err := svc.PutItem(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

func BulkInsert(ctx context.Context, items []SerializeAble) error {
	svc := ctx.Value("dynamodb").(*dynamodb.Client)
	var requests []types.WriteRequest
	for _, item := range items {
		av := item.Deserialize()
		requests = append(requests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		})
	}
	tableName := items[0].TableName(ctx)
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: requests,
		},
	}
	_, err := svc.BatchWriteItem(context.Background(), input)

	if err != nil {
		return err
	}
	return nil
}

func Query(ctx context.Context, tableName string, query *InputQuery) ([]map[string]types.AttributeValue, error) {
	svc := ctx.Value("dynamodb").(*dynamodb.Client)
	input := &dynamodb.QueryInput{
		TableName:        aws.String(tableName),
		ScanIndexForward: aws.Bool(query.ScanIndexForward),
	}
	if query.ExpressionAttribute != nil {
		input.ExpressionAttributeValues = *(query.ExpressionAttribute)
	}
	if query.KeyConditionExpression != nil {
		input.KeyConditionExpression = aws.String(*(query.KeyConditionExpression))
	}
	if query.FilterExpression != nil {
		input.FilterExpression = aws.String(*(query.FilterExpression))
	}
	if query.Limit != 0 && query.Limit <= 100 {
		input.Limit = aws.Int32(query.Limit)
	} else {
		input.Limit = aws.Int32(100)
	}
	if query.IndexName != "" {
		input.IndexName = aws.String(query.IndexName)
	}
	if query.ExclusiveStartKey != nil {
		input.ExclusiveStartKey = *(query.ExclusiveStartKey)
	}
	output, err := svc.Query(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return output.Items, nil
}

func Update(ctx context.Context, tableName string, query *UpdateQuery) error {
	svc := ctx.Value("dynamodb").(*dynamodb.Client)
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
	}
	if query.ExpressionAttribute != nil {
		input.ExpressionAttributeValues = *(query.ExpressionAttribute)
	}
	if query.UpdateExpression != nil {
		input.UpdateExpression = aws.String(*(query.UpdateExpression))
	}
	if query.Key != nil {
		input.Key = *(query.Key)
	}
	if query.ReturnValues != "" {
		input.ReturnValues = query.ReturnValues
	}

	_, err := svc.UpdateItem(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
