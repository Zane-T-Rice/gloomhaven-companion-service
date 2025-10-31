package main

import (
	"context"
	setenvironmentvariables "gloomhaven-companion-service/internal/set-environment-variables"
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

type Item struct {
	Id     string `dynamodbav:"id"`
	Entity string `dynamodbav:"entity"`
}

var dynamoDbClient *dynamodb.Client

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Handling websocket disconnnect request")

	setenvironmentvariables.SetEnvironmentVariables()
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

	scenarioId := request.Body
	log.Printf("scenarioId=%s", scenarioId)

	tableName := "gloomhaven-companion-service"
	connections, err := dynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		ProjectionExpression:   aws.String("id, entity"),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: request.RequestContext.ConnectionID},
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
				"id":     item["id"],
				"entity": item["entity"],
			},
		})
	}
	if err != nil {
		log.Printf("Couldn't delete item from table. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(handleRequest)
}
