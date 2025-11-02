package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterCampaignsRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	campaignsController := controllers.NewCampaignsController(dynamodb)

	campaign := app.Group("/" + constants.CAMPAIGNS)
	campaign.Get("/",
		campaignsController.List,
	)
	campaign.Post("/",
		campaignsController.Create,
	)
	campaign.Patch("/:campaignId",
		campaignsController.Patch,
	)
	campaign.Delete("/:campaignId",
		campaignsController.Delete,
	)
}
