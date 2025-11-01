package utils

import (
	"context"
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	DynamoDBClient *dynamodb.Client
}

func (db *DynamoDB) ConnectToDynamoDB() {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv(constants.LOCAL_SERVICE_PORT) == "" {
		db.DynamoDBClient = dynamodb.NewFromConfig(config)
	} else {
		db.DynamoDBClient = dynamodb.NewFromConfig(config, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(os.Getenv(constants.LOCAL_DATABASE_ENDPOINT))
		})
	}
}

func (db *DynamoDB) PutItem(item interface{}) error {
	dynamodbItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	_, err = db.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(constants.TABLE_NAME), Item: dynamodbItem,
	})
	if err != nil {
		log.Fatalf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return nil
}

func buildUpdateExpression(input interface{}) (*expression.Expression, error) {
	inputItem, err := attributevalue.MarshalMap(input)
	if err != nil {
		return nil, err
	}
	var updateBuilder expression.UpdateBuilder
	for key, val := range inputItem {
		updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(val))
	}
	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	return &expr, err
}

func (db *DynamoDB) UpdateItem(
	partitionKey string,
	partitionKeyValue string,
	sortKey string,
	sortKeyValue string,
	input interface{},
	output interface{},
) error {
	expr, err := buildUpdateExpression(input)
	if err != nil {
		return err
	}
	updateItemResult, err := db.DynamoDBClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(constants.TABLE_NAME),
		Key: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{Value: partitionKeyValue},
			sortKey:      &types.AttributeValueMemberS{Value: sortKeyValue},
		},
		ReturnValues:              types.ReturnValueAllNew,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
	})
	if err != nil {
		log.Fatalf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	err = attributevalue.UnmarshalMap(updateItemResult.Attributes, &output)
	if err != nil {
		log.Fatalf("failed to unmarshal DynamoDB item, %v", err)
		return err
	}
	return nil
}

func (db *DynamoDB) GetItem(
	partitionKey string,
	partitionKeyValue string,
	sortKey string,
	sortKeyValue string,
	output interface{},
) error {
	result, err := db.DynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(constants.TABLE_NAME),
		Key: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{Value: partitionKeyValue},
			sortKey:      &types.AttributeValueMemberS{Value: sortKeyValue},
		},
	})
	if err != nil {
		log.Fatalf("Couldn't get item from table. Here's why: %v\n", err)
		return err
	}
	if result.Item == nil {
		return errors.NewNotFoundError()
	}
	err = attributevalue.UnmarshalMap(result.Item, &output)
	if err != nil {
		log.Fatalf("failed to unmarshal DynamoDB item, %v", err)
		return err
	}
	return nil
}

func (db *DynamoDB) Query(
	partitionKey string,
	partitionKeyValue string,
	sortKey string,
	sortKeyValue string,
	indexName *string,
	output any, // output must be a non-nil pointer
) error {
	queryInput := dynamodb.QueryInput{
		TableName: aws.String(constants.TABLE_NAME),
		KeyConditions: map[string]types.Condition{
			partitionKey: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: partitionKeyValue},
				},
			},
			sortKey: {
				ComparisonOperator: types.ComparisonOperatorBeginsWith,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: sortKeyValue},
				},
			},
		},
		IndexName: indexName,
	}
	result, err := db.DynamoDBClient.Query(context.TODO(), &queryInput)
	if err != nil {
		log.Fatalf("Couldn't get item from table. Here's why: %v\n", err)
		return err
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, output)
	if err != nil {
		log.Fatalf("failed to unmarshal DynamoDB items, %v", err)
		return err
	}

	return nil
}

func (db *DynamoDB) DeleteItem(
	partitionKey string,
	partitionKeyValue string,
	sortKey string,
	sortKeyValue string,
	output any, // output must be a non-nil pointer
) error {
	deleteItemResults, err := db.DynamoDBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(constants.TABLE_NAME),
		Key: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{Value: partitionKeyValue},
			sortKey:      &types.AttributeValueMemberS{Value: sortKeyValue},
		},
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		log.Fatalf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	err = attributevalue.UnmarshalMap(deleteItemResults.Attributes, &output)
	if err != nil {
		log.Fatalf("failed to unmarshal DynamoDB item, %v", err)
		return err
	}
	return nil
}

func NewDynamoDB() DynamoDB {
	return DynamoDB{}
}
