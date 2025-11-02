package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterScenariosRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	scenariosController := controllers.NewScenariosController(dynamodb)

	scenario := app.Group("/" + constants.CAMPAIGNS + "/:campaignId/" + constants.SCENARIOS)
	scenario.Get("/",
		scenariosController.List,
	)
	scenario.Post("/",
		scenariosController.Create,
	)
	scenario.Patch("/:scenarioId",
		scenariosController.Patch,
	)
	scenario.Delete("/:scenarioId",
		scenariosController.Delete,
	)
}
