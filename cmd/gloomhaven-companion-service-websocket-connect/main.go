package main

import (
	"context"
	"gloomhaven-companion-service/internal/websockets"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	return websockets.Connect(ctx, websockets.ConnectRequest{
		Headers:               request.Headers,
		QueryStringParameters: request.QueryStringParameters,
		RequestContext: websockets.RequestContext{
			ConnectionID: request.RequestContext.ConnectionID,
		},
	})
}

func main() {
	lambda.Start(handleRequest)
}
