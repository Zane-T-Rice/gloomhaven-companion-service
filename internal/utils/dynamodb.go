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
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
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
			"parent": &types.AttributeValueMemberS{Value: partitionKeyValue},
			"entity": &types.AttributeValueMemberS{Value: sortKeyValue},
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
		log.Fatalf("failed to unmarshal DynamoDB item, %v", err)
		return err
	}
	return nil
}

func NewDynamoDB() DynamoDB {
	return DynamoDB{}
}
