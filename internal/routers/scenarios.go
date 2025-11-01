package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	middlewares "gloomhaven-companion-service/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterScenariosRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	scenariosController := controllers.NewScenariosController(dynamodb)

	scenario := app.Group("/" + constants.CAMPAIGNS + "/:campaignId/" + constants.SCENARIOS)
	scenario.Get("/",
		middlewares.HasScope(constants.SCOPE_READ_CAMPAIGNS),
		scenariosController.List,
	)
	scenario.Post("/",
		middlewares.HasScope(constants.SCOPE_WRITE_CAMPAIGNS),
		scenariosController.Create,
	)
	scenario.Patch("/:scenarioId",
		middlewares.HasScope(constants.SCOPE_WRITE_CAMPAIGNS),
		scenariosController.Patch,
	)
	scenario.Delete("/:scenarioId",
		middlewares.HasScope(constants.SCOPE_WRITE_CAMPAIGNS),
		scenariosController.Delete,
	)
}
