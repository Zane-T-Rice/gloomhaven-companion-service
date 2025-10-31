// main.go
package main

import (
	"context"
	"gloomhaven-companion-service/internal/middlewares"
	"gloomhaven-companion-service/internal/routers"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2" // or v3
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"

	utils "gloomhaven-companion-service/internal/utils"
)

var app *fiber.App
var fiberLambda *fiberadapter.FiberLambda

// init the Fiber Server
func init() {
	utils.SetEnvironmentVariables()
	dynamodb := utils.NewDynamoDB()
	dynamodb.ConnectToDynamoDB()

	app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("WEBSITE_DOMAIN"),
	}))

	// Always ensure the token is valid before doing anything.
	app.Use(adaptor.HTTPMiddleware(middlewares.EnsureValidToken()))

	routers.RegisterCampaignsRoutes(app, dynamodb)

	fiberLambda = fiberadapter.New(app)
}

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// I cannot figure out how to prevent API Gateway from including the base path
	// in the request path, so we manually strip it here.
	req.Path, _ = strings.CutPrefix(req.Path, "/gloomhaven-companion-service")
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	localServicePort := os.Getenv("LOCAL_SERVICE_PORT")
	if localServicePort == "" {
		lambda.Start(Handler)
	} else {
		// Run the service locally without the Lambda wrapper.
		log.Fatal(app.Listen(":" + localServicePort))
	}
}
