package utils

import (
	"context"
	"gloomhaven-companion-service/internal/constants"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func NewDynamoDB() DynamoDB {
	return DynamoDB{}
}
