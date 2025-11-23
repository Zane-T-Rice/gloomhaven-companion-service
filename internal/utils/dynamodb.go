package utils

import (
	"context"
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/errors"
	"log"
	"os"
	"time"

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
		log.Print(err)
	}

	if os.Getenv(constants.LOCAL_SERVICE_PORT) == "" {
		db.DynamoDBClient = dynamodb.NewFromConfig(config)
	} else {
		db.DynamoDBClient = dynamodb.NewFromConfig(config, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(os.Getenv(constants.LOCAL_DATABASE_ENDPOINT))
		})
	}
}

func (db *DynamoDB) PutItem(item any) error {
	dynamodbItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	_, err = db.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(constants.TABLE_NAME), Item: dynamodbItem,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return nil
}

func buildUpdateExpression(input any) (*expression.Expression, error) {
	inputItem, err := attributevalue.MarshalMap(input)
	if err != nil {
		return nil, err
	}
	// Set condition that the incoming updated_at must match the database.
	conditionBuilder := expression.Equal(expression.Name("updated_at"), expression.Value(inputItem["updated_at"]))
	// Set the new updated_at that will be written if the above condition is met.
	inputItem["updated_at"] = &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)}
	var updateBuilder expression.UpdateBuilder
	for key, val := range inputItem {
		updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(val))
	}
	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).WithCondition(conditionBuilder).Build()
	return &expr, err
}

func (db *DynamoDB) UpdateItem(
	partitionKey string,
	partitionKeyValue string,
	sortKey string,
	sortKeyValue string,
	input any,
	output any,
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
		log.Printf("Couldn't update item in table. Here's why: %v\n", err)
		return err
	}
	err = attributevalue.UnmarshalMap(updateItemResult.Attributes, &output)
	if err != nil {
		log.Printf("failed to unmarshal DynamoDB item, %v", err)
		return err
	}
	return nil
}

func (db *DynamoDB) GetItem(
	partitionKey string,
	partitionKeyValue string,
	sortKey string,
	sortKeyValue string,
	output any,
) error {
	result, err := db.DynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(constants.TABLE_NAME),
		Key: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{Value: partitionKeyValue},
			sortKey:      &types.AttributeValueMemberS{Value: sortKeyValue},
		},
	})
	if err != nil {
		log.Printf("Couldn't get item from table. Here's why: %v\n", err)
		return err
	}
	if result.Item == nil {
		return errors.NewNotFoundError()
	}
	err = attributevalue.UnmarshalMap(result.Item, &output)
	if err != nil {
		log.Printf("failed to unmarshal DynamoDB item, %v", err)
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
		log.Printf("Couldn't get item from table. Here's why: %v\n", err)
		return err
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, output)
	if err != nil {
		log.Printf("failed to unmarshal DynamoDB items, %v", err)
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
		log.Printf("Couldn't delete item from table. Here's why: %v\n", err)
		return err
	}
	err = attributevalue.UnmarshalMap(deleteItemResults.Attributes, &output)
	if err != nil {
		log.Printf("failed to unmarshal DynamoDB item, %v", err)
		return err
	}
	return nil
}

func NewDynamoDB() DynamoDB {
	return DynamoDB{}
}
