package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"

	middlewares "gloomhaven-companion-service/internal/middlewares"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
)

func RegisterCampaignsRoutes(app *fiber.App, dynamodbClient *dynamodb.Client) {
	campaignsController := controllers.CampaignsController(dynamodbClient)

	campaign := app.Group("/campaigns")
	campaign.Get("/",
		middlewares.HasScope(constants.SCOPE_READ_CAMPAIGNS),
		campaignsController.List,
	)
}
