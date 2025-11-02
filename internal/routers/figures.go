package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterFiguresRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	figuresController := controllers.NewFiguresController(dynamodb)

	scenario := app.Group("/" + constants.CAMPAIGNS + "/:campaignId/" + constants.SCENARIOS + "/:scenarioId/" + constants.FIGURES)
	scenario.Get("/",
		figuresController.List,
	)
	scenario.Post("/",
		figuresController.Create,
	)
	scenario.Patch("/:figureId",
		figuresController.Patch,
	)
	scenario.Delete("/:figureId",
		figuresController.Delete,
	)
}
