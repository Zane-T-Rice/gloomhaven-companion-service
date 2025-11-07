package websockets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/websocket"
)

type DefaultRequest struct {
	Body           string
	RequestContext RequestContext
}

func Default(ctx context.Context, request DefaultRequest) (events.APIGatewayProxyResponse, error) {
	var dynamoDbClient *dynamodb.Client

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

	figure := types.FigureItem{}
	json.Unmarshal([]byte(request.Body), &figure)

	tableName := "gloomhaven-companion-service"

	connections, err := dynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		ProjectionExpression:   aws.String("parent, entity"),
		KeyConditionExpression: aws.String("entity = :entity"),
		ExpressionAttributeValues: map[string]awsTypes.AttributeValue{
			":entity": &awsTypes.AttributeValueMemberS{Value: figure.Parent},
		},
		IndexName: aws.String("entity-index"),
	})
	if err != nil {
		log.Printf("Couldn't query items from table. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Broadcast the figure object to all the listeners except the sender.
	svc := apigatewaymanagementapi.NewFromConfig(config, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String(os.Getenv(constants.WEB_SOCKETS_URL))
	})
	for _, item := range connections.Items {
		connectionId := aws.String(item["parent"].(*awsTypes.AttributeValueMemberS).Value)
		if *connectionId == request.RequestContext.ConnectionID {
			continue
		}
		err = nil
		if os.Getenv("LOCAL_SERVICE_PORT") == "" {
			_, err = svc.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: connectionId,
				Data:         []byte(request.Body),
			})
		} else {
			connection := Connections[item["parent"].(*awsTypes.AttributeValueMemberS).Value]
			if connection == nil {
				err = errors.New("connection does not exist")
			}
			if err == nil {
				if err = connection.WriteMessage(websocket.TextMessage, []byte(request.Body)); err != nil {
					fmt.Println("Error writing message:", err)
				}
			}
		}
		if err != nil {
			log.Printf("Error posting to connection %s: %v", *connectionId, err)

			// Delete the bad connection.
			dynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key: map[string]awsTypes.AttributeValue{
					"parent": item["parent"],
					"entity": item["entity"],
				},
			})
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}
