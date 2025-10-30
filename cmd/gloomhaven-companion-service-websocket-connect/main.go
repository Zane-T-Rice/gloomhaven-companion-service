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

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
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
		return generatePolicy("user", "Deny", "*"), nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", os.Getenv("GLOOMHAVEN_COMPANION_SERVICE_URL")+"/enemies", nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return generatePolicy("user", "Deny", "*"), nil
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error calling companion service: %v", err)
		return generatePolicy("user", "Deny", "*"), nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("companion service response: status=%d body=%s", resp.StatusCode, string(body))

	// Treat HTTP 200 as a successful validation. Adjust logic if the companion
	// service uses a different success status or response shape.
	if resp.StatusCode == http.StatusOK {
		// TODO: put the connectionId and scenario into DynamoDB here if needed.
		return generatePolicy("user", "Allow", "*"), nil
	}

	// Anything else is a denial
	log.Printf("authorization denied by companion service: status=%d", resp.StatusCode)
	return generatePolicy("user", "Deny", "*"), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		statement := events.IAMPolicyStatement{
			Action:   []string{"execute-api:Invoke"},
			Effect:   effect,
			Resource: []string{resource},
		}
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				statement,
			},
		}
	}
	return authResponse
}

func main() {
	lambda.Start(handleRequest)
}
