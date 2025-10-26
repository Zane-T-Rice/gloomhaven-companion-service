// main.go
package main

import (
	"context"
	"encoding/json"
	ensurevalidtoken "gloomhaven-companion-service/internal/ensure-valid-token"
	"log"
	"maps"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor" // or v3
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var app *fiber.App
var fiberLambda *fiberadapter.FiberLambda

// init the Fiber Server
func init() {
	godotenv.Load()
	localServicePort := os.Getenv("LOCAL_SERVICE_PORT")
	if localServicePort == "" {
		// Load any environment variables from AWS Secret Manager
		audienceSecretName := "gloomhaven-companion-service-audience"
		issuerSecretName := "gloomhaven-companion-service-issuer"
		region := "us-east-1"

		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatal(err)
		}

		// Create Secrets Manager client
		svc := secretsmanager.NewFromConfig(config)

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(audienceSecretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		}
		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			log.Fatal(err.Error())
		}
		type Secret struct {
			Secret string `json:"value"`
		}
		audienceSecret := Secret{}
		json.Unmarshal([]byte(*result.SecretString), &audienceSecret)

		input = &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(issuerSecretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		}
		result, err = svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			log.Fatal(err.Error())
		}

		issuerSecret := Secret{}
		json.Unmarshal([]byte(*result.SecretString), &issuerSecret)

		log.Println("Setting the AUDIENCE and ISSUER")
		os.Setenv("AUDIENCE", audienceSecret.Secret)
		os.Setenv("ISSUER", issuerSecret.Secret)
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

// I need this hideous function becuase fiberLambda only works with
// APIGatewayProxyRequest/APIGatewayProxyResponse and I am using
// a Lambda function url at the moment.
func convertLambdaFunctionURLRequestToAPIGatewayProxyRequest(
	lfurlReq events.LambdaFunctionURLRequest,
) events.APIGatewayProxyRequest {
	apiGatewayReq := events.APIGatewayProxyRequest{
		Path:            lfurlReq.RawPath,
		HTTPMethod:      lfurlReq.RequestContext.HTTP.Method,
		Headers:         lfurlReq.Headers,
		Body:            lfurlReq.Body,
		IsBase64Encoded: lfurlReq.IsBase64Encoded,
	}

	// Add query string parameters
	apiGatewayReq.QueryStringParameters = make(map[string]string)
	maps.Copy(apiGatewayReq.QueryStringParameters, lfurlReq.QueryStringParameters)

	return apiGatewayReq
}

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.APIGatewayProxyResponse, error) {
	apiGatewayProxyRequest := convertLambdaFunctionURLRequestToAPIGatewayProxyRequest(req)
	return fiberLambda.ProxyWithContext(ctx, apiGatewayProxyRequest)
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
