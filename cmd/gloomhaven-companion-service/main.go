// main.go
package main

import (
	"context"
	"gloomhaven-companion-service/internal/constants"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor" // or v3
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	middleware "gloomhaven-companion-service/internal/middleware"
	utils "gloomhaven-companion-service/internal/utils"
)

var app *fiber.App
var fiberLambda *fiberadapter.FiberLambda
var dynamoDbClient *dynamodb.Client

type Item struct {
	Id     string `dynamodbav:"id"`
	Entity string `dynamodbav:"entity"`
}

// init the Fiber Server
func init() {
	utils.SetEnvironmentVariables()
	utils.ConnectToDynamoDB(&dynamoDbClient)

	app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("WEBSITE_DOMAIN"),
	}))

	app.Get("/campaign",
		adaptor.HTTPMiddleware(middleware.EnsureValidToken()),
		middleware.HasScope(constants.SCOPE_READ_ENEMIES),
		func(c *fiber.Ctx) error {
			// itemDTO := Item{
			// 	Id:     "enemy1",
			// 	Entity: "#ENEMY",
			// }

			// tableName := "gloomhaven-companion-service"

			// item, err := attributevalue.MarshalMap(itemDTO)
			// if err != nil {
			// 	panic(err)
			// }
			// _, err = dynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			// 	TableName: aws.String(tableName), Item: item,
			// })
			// if err != nil {
			// 	log.Printf("Couldn't add item to table. Here's why: %v\n", err)
			// 	return err
			// }
			return c.SendString("[List of enemies]")
		},
	)

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
