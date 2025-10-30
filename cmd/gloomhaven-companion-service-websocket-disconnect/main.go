package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) error {
	log.Printf("HANDLING THE DISCONNECT ROUTE REQUEST: %+v", request)
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
