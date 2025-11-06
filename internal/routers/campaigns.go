package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterCampaignsRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	campaignsController := controllers.NewCampaignsController(dynamodb)

	// This lets people use their Join Codes to join a Campaign.
	app.Post("/join",
		campaignsController.JoinCampaign,
	)

	campaign := app.Group("/" + constants.CAMPAIGNS)

	// Campaign CRUD
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

	// This lets someone who can see a campaign generate an invite code for it.
	campaign.Post("/:campaignId/create-join-code",
		campaignsController.CreateJoinCode,
	)
}
