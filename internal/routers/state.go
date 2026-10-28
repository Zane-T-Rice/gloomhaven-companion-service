package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterStateRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	stateController := controllers.NewStateController(dynamodb)

	state := app.Group("/" + constants.CAMPAIGNS + "/:campaignId/" + constants.SCENARIOS + "/:scenarioId/" + constants.STATE_PATH)
	state.Get("/",
		stateController.List,
	)
	state.Put("/",
		stateController.Put,
	)
}
