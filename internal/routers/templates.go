package routers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Any player can make, change and delete their own templates and editions. Admins can also do that
// to the system's, which are asked for with ?scope=system when creating. The services decide who
// may change what.
func RegisterTemplatesRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	templatesController := controllers.NewTemplatesController(dynamodb)

	// These come before /templates/:templateId, so that "editions" is not taken for a template id.
	editions := app.Group("/" + constants.TEMPLATES + "/" + constants.EDITIONS)
	editions.Get("/", templatesController.ListEditions)
	editions.Post("/", templatesController.CreateEdition)
	editions.Patch("/:editionId", templatesController.PatchEdition)
	editions.Delete("/:editionId", templatesController.DeleteEdition)

	templates := app.Group("/" + constants.TEMPLATES)
	templates.Get("/", templatesController.List)
	templates.Post("/", templatesController.Create)
	templates.Patch("/:templateId", templatesController.Patch)
	templates.Delete("/:templateId", templatesController.Delete)
}
