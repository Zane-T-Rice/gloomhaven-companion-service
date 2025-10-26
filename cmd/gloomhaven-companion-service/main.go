// main.go
package main

import (
	"context"
	ensurevalidtoken "gloomhaven-companion-service/internal/ensure-valid-token"
	"log"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	// Load any environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	app = fiber.New()

	app.Get("/enemies",
		adaptor.HTTPMiddleware(ensurevalidtoken.EnsureValidToken()),
		func(c *fiber.Ctx) error {
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
			return c.SendString("Hello, World 4!")
		},
	)

	fiberLambda = fiberadapter.New(app)
}

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
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
