package controllers

import (
	"gloomhaven-companion-service/internal/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func NewCampaignsController(dynamodbClient *dynamodb.Client) CampaignsController {
	campaignsService := services.NewCampaignsService(dynamodbClient)

	return CampaignsController{
		CampaignsService: campaignsService,
	}
}
