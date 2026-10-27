package controllers

import "github.com/gofiber/fiber/v2"

// Pinger is anything that can tell whether the database can be reached.
type Pinger interface {
	Ping() error
}

type HealthController struct {
	Database Pinger
}

func NewHealthController(database Pinger) HealthController {
	return HealthController{Database: database}
}

// Check answers 200 when the service can read its table and 503 when it cannot.
func (c *HealthController) Check(cxt *fiber.Ctx) error {
	if err := c.Database.Ping(); err != nil {
		return cxt.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
	}
	return cxt.JSON(fiber.Map{"status": "ok"})
}
