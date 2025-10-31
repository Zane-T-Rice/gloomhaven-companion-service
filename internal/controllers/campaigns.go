package controllers

import (
	"gloomhaven-companion-service/internal/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
)

type campaignsController struct {
	List func(c *fiber.Ctx) error
}

func CampaignsController(dynamodbClient *dynamodb.Client) campaignsController {
	campaignsService := services.CampaignsService(dynamodbClient)

	return campaignsController{
		List: func(c *fiber.Ctx) error {
			campaigns, err := campaignsService.List()
			if err != nil {
				return err
			}
			return c.JSON(campaigns)
		},
	}
}
