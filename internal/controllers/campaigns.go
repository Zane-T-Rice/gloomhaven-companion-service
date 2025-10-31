package controllers

import (
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CampaignsController struct {
	CampaignsService services.CampaignsService
}

func (c CampaignsController) List(cxt *fiber.Ctx) error {
	campaigns, err := c.CampaignsService.List()
	if err != nil {
		return err
	}
	return cxt.JSON(campaigns)
}

func NewCampaignsController(dynamodb utils.DynamoDB) CampaignsController {
	campaignsService := services.NewCampaignsService(dynamodb)

	return CampaignsController{
		CampaignsService: campaignsService,
	}
}
