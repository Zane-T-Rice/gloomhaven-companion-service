package websockets

import (
	"context"
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type RequestContext struct {
	ConnectionID string
}
type ConnectRequest struct {
	Headers               map[string]string
	QueryStringParameters map[string]string
	RequestContext        RequestContext
}

func Connect(ctx context.Context, request ConnectRequest) (events.APIGatewayProxyResponse, error) {
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

	// Extract the authorization token and scenarioId from the request
	protocol := request.Headers["Sec-WebSocket-Protocol"]
	token := "Bearer " + protocol
	campaignId := request.QueryStringParameters["campaignId"]
	scenarioId := request.QueryStringParameters["scenarioId"]

	// Make a GET call to the companion service to validate the token.
	// The companion service expects the Authorization header to be forwarded.
	if token == "" {
		log.Println("no Authorization header provided")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		os.Getenv("GLOOMHAVEN_COMPANION_SERVICE_URL")+
			"/campaigns/"+campaignId+"/scenarios/"+scenarioId+"/figures", nil,
	)
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

	if resp.StatusCode == http.StatusOK {
		itemDTO := types.Item{
			Parent: request.RequestContext.ConnectionID,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId +
				constants.SCENARIO + constants.SEPERATOR + scenarioId,
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
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: map[string]string{
			"Sec-WebSocket-Protocol": protocol,
		}}, nil
	}

	// Anything else is a denial
	log.Printf("authorization denied by companion service: status=%d", resp.StatusCode)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
}
