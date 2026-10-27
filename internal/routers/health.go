package routers

import (
	"gloomhaven-companion-service/internal/controllers"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(app *fiber.App, dynamodb utils.DynamoDB) {
	healthController := controllers.NewHealthController(&dynamodb)

	app.Get("/health", healthController.Check)
}
