package main

import (
	"context"
	"gloomhaven-companion-service/internal/utils"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamoDbClient *dynamodb.Client

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	utils.SetEnvironmentVariables()
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("LOCAL_SERVICE_PORT") == "" {
		dynamoDbClient = dynamodb.NewFromConfig(config)
	} else {
		dynamoDbClient = dynamodb.NewFromConfig(config, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://localhost:8000/")
		})
	}

	tableName := "gloomhaven-companion-service"
	connections, err := dynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		ProjectionExpression:   aws.String("parent, entity"),
		KeyConditionExpression: aws.String("parent = :parent"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":parent": &types.AttributeValueMemberS{Value: request.RequestContext.ConnectionID},
		},
	})
	if err != nil {
		log.Printf("Couldn't get the item from table to delete it. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	// There should only be one, but any record with this connection id should be deleted.
	for _, item := range connections.Items {
		_, err = dynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"parent": item["parent"],
				"entity": item["entity"],
			},
		})
		if err != nil {
			log.Printf("Couldn't delete item from table. Here's why: %v\n", err)
		} else {
			log.Printf("Deleted connection id %s", request.RequestContext.ConnectionID)
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(handleRequest)
}
