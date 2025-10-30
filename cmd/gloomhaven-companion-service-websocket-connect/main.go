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
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Handling websocket connect request")

	setenvironmentvariables.SetEnvironmentVariables()

	// Extract the authorization token and scenarioId from the request
	token := request.Headers["Authorization"]
	scenario := request.QueryStringParameters["scenarioId"]
	log.Printf("scenarioId=%s", scenario)

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
		// TODO: put the connectionId and scenario into DynamoDB here if needed.
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	}

	// Anything else is a denial
	log.Printf("authorization denied by companion service: status=%d", resp.StatusCode)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
}

func main() {
	lambda.Start(handleRequest)
}
