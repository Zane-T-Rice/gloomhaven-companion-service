package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/middlewares"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterTemplatesRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	templatesController := controllers.NewTemplatesController(dynamodb)

	scenario := app.Group("/" + constants.TEMPLATES)
	scenario.Get("/",
		templatesController.List,
	)
	scenario.Post("/",
		middlewares.HasOneOfScopes([]string{constants.SCOPE_ADMIN}),
		templatesController.Create,
	)
	scenario.Patch("/:templateId",
		middlewares.HasOneOfScopes([]string{constants.SCOPE_ADMIN}),
		templatesController.Patch,
	)
	scenario.Delete("/:templateId",
		middlewares.HasOneOfScopes([]string{constants.SCOPE_ADMIN}),
		templatesController.Delete,
	)
}
