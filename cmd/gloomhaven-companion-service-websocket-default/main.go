package main

import (
	"context"
	"gloomhaven-companion-service/internal/websockets"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	return websockets.Default(ctx, websockets.DefaultRequest{
		Body: request.Body,
		RequestContext: websockets.RequestContext{
			ConnectionID: request.RequestContext.ConnectionID,
		},
	})
}

func main() {
	lambda.Start(handleRequest)
}
