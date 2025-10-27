// main.go
package main

import (
	"context"
	ensurevalidtoken "gloomhaven-companion-service/internal/ensure-valid-token"
	setenvironmentvariables "gloomhaven-companion-service/internal/set-environment-variables"
	"log"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor" // or v3
	"github.com/valyala/fasthttp"
)

var app *fiber.App
var fiberLambda *fiberadapter.FiberLambda

// init the Fiber Server
func init() {
	setenvironmentvariables.SetEnvironmentVariables()

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
// func convertLambdaFunctionURLRequestToAPIGatewayProxyRequest(
// 	lfurlReq events.LambdaFunctionURLRequest,
// ) events.APIGatewayProxyRequest {
// 	apiGatewayReq := events.APIGatewayProxyRequest{
// 		Path:            lfurlReq.RawPath,
// 		HTTPMethod:      lfurlReq.RequestContext.HTTP.Method,
// 		Headers:         lfurlReq.Headers,
// 		Body:            lfurlReq.Body,
// 		IsBase64Encoded: lfurlReq.IsBase64Encoded,
// 	}
//
// 	// Add query string parameters
// 	apiGatewayReq.QueryStringParameters = make(map[string]string)
// 	maps.Copy(apiGatewayReq.QueryStringParameters, lfurlReq.QueryStringParameters)
//
// 	return apiGatewayReq
// }

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
