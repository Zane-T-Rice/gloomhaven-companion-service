package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	middlewares "gloomhaven-companion-service/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterFiguresRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	figuresController := controllers.NewFiguresController(dynamodb)

	scenario := app.Group("/" + constants.CAMPAIGNS + "/:campaignId/" + constants.SCENARIOS + "/:scenarioId/" + constants.FIGURES)
	scenario.Get("/",
		middlewares.HasScope(constants.SCOPE_READ_CAMPAIGNS),
		figuresController.List,
	)
	scenario.Post("/",
		middlewares.HasScope(constants.SCOPE_WRITE_CAMPAIGNS),
		figuresController.Create,
	)
	scenario.Patch("/:figureId",
		middlewares.HasScope(constants.SCOPE_WRITE_CAMPAIGNS),
		figuresController.Patch,
	)
	scenario.Delete("/:figureId",
		middlewares.HasScope(constants.SCOPE_WRITE_CAMPAIGNS),
		figuresController.Delete,
	)
}
