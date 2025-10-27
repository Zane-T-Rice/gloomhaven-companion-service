// main.go
package main

import (
	"context"
	ensurevalidtoken "gloomhaven-companion-service/internal/ensure-valid-token"
	setenvironmentvariables "gloomhaven-companion-service/internal/set-environment-variables"
	"log"
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor" // or v3
	"github.com/valyala/fasthttp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	app = fiber.New()

	app.Get("/enemies",
		adaptor.HTTPMiddleware(ensurevalidtoken.EnsureValidToken()),
		func(c *fiber.Ctx) error {
			log.Println("Starting /enemies")
			token := c.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

			claims := token.CustomClaims.(*ensurevalidtoken.CustomClaims)
			if !claims.HasScope("read:enemies") {
				c.Response().SetStatusCode(fasthttp.StatusForbidden)
				c.Write([]byte(`{"message":"Insufficient scope."}`))
				return nil
			}

			return c.Next()
		},
		func(c *fiber.Ctx) error {
			itemDTO := Item{
				Id:     "enemy1",
				Entity: "#ENEMY",
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
				return err
			}
			return c.SendString("[List of enemies]")
		},
	)

	app.Get("/characters",
		adaptor.HTTPMiddleware(ensurevalidtoken.EnsureValidToken()),
		func(c *fiber.Ctx) error {
			log.Println("Starting /characters")
			token := c.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

			claims := token.CustomClaims.(*ensurevalidtoken.CustomClaims)
			if !claims.HasScope("read:characters") {
				c.Response().SetStatusCode(fasthttp.StatusForbidden)
				c.Write([]byte(`{"message":"Insufficient scope."}`))
				return nil
			}

			return c.Next()
		},
		func(c *fiber.Ctx) error {
			return c.SendString("[List of characters]")
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
