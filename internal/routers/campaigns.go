package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	middlewares "gloomhaven-companion-service/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterCampaignsRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	campaignsController := controllers.NewCampaignsController(dynamodb)

	campaign := app.Group("/" + constants.CAMPAIGNS)
	campaign.Get("/",
		middlewares.HasScope(constants.SCOPE_READ_CAMPAIGNS),
		campaignsController.List,
	)
}
