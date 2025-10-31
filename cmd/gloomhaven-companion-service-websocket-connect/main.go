package main

import (
	"context"
	setenvironmentvariables "gloomhaven-companion-service/internal/set-environment-variables"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Item struct {
	Id     string `dynamodbav:"id"`
	Entity string `dynamodbav:"entity"`
}

var dynamoDbClient *dynamodb.Client

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Handling websocket connect request")

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

	// Extract the authorization token and scenarioId from the request
	token := request.Headers["Authorization"]
	scenarioId := request.QueryStringParameters["scenarioId"]
	log.Printf("scenarioId=%s", scenarioId)

	// Make a GET call to the companion service to validate the token.
	// The companion service expects the Authorization header to be forwarded.
	if token == "" {
		log.Println("no Authorization header provided")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", os.Getenv("GLOOMHAVEN_COMPANION_SERVICE_URL")+"/enemies", nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error calling companion service: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("companion service response: status=%d body=%s", resp.StatusCode, string(body))

	// Treat HTTP 200 as a successful validation. Adjust logic if the companion
	// service uses a different success status or response shape.
	if resp.StatusCode == http.StatusOK {
		itemDTO := Item{
			Id:     request.RequestContext.ConnectionID,
			Entity: scenarioId,
		}

		tableName := "gloomhaven-companion-service"

		item, err := attributevalue.MarshalMap(itemDTO)
		if err != nil {
			panic(err)
		}
		_, err = dynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName), Item: item,
		})
		if err != nil {
			log.Printf("Couldn't add item to table. Here's why: %v\n", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	}

	// Anything else is a denial
	log.Printf("authorization denied by companion service: status=%d", resp.StatusCode)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
}

func main() {
	lambda.Start(handleRequest)
}
