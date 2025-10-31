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
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Item struct {
	Id     string `dynamodbav:"id"`
	Entity string `dynamodbav:"entity"`
}

var dynamoDbClient *dynamodb.Client

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Handling websocket default request")

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
		KeyConditionExpression: aws.String("entity = :entity"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":entity": &types.AttributeValueMemberS{Value: scenarioId},
		},
		IndexName: aws.String("entity-index"),
	})
	if err != nil {
		log.Printf("Couldn't query items from table. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	svc := apigatewaymanagementapi.NewFromConfig(config, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String("https://ws.zanesworld.click/gloomhaven-companion-service")
	})

	for _, item := range connections.Items {
		connectionId := aws.String(item["id"].(*types.AttributeValueMemberS).Value)
		_, err := svc.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: connectionId,
			Data:         []byte(scenarioId),
		})
		if err != nil {
			log.Printf("Error posting to connection %s: %v", connectionId, err)
			// Handle errors, e.g., remove invalid connection ID from store
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(handleRequest)
}
