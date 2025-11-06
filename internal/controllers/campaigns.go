package controllers

import (
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"
)

type CampaignsController struct {
	CampaignsService services.CampaignsService
}

func (c *CampaignsController) List(cxt *fiber.Ctx) error {
	token := cxt.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	playerId := token.RegisteredClaims.Subject
	campaigns, err := c.CampaignsService.List(playerId)
	if err != nil {
		return err
	}
	return cxt.JSON(campaigns)
}

func (c *CampaignsController) Create(cxt *fiber.Ctx) error {
	input := types.CampaignCreateInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	token := cxt.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	playerId := token.RegisteredClaims.Subject
	campaign, err := c.CampaignsService.Create(input, playerId)
	if err != nil {
		return err
	}
	return cxt.JSON(*campaign)
}

func (c *CampaignsController) Patch(cxt *fiber.Ctx) error {
	input := types.CampaignPatchInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	campaignId := cxt.Params("campaignId")
	campaign, err := c.CampaignsService.Patch(input, campaignId)
	if err != nil {
		return err
	}
	return cxt.JSON(*campaign)
}

func (c *CampaignsController) Delete(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	campaign, err := c.CampaignsService.Delete(campaignId)
	if err != nil {
		return err
	}
	return cxt.JSON(*campaign)
}

func (c *CampaignsController) CreateJoinCode(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	joinCampaignCode, err := c.CampaignsService.CreateJoinCode(campaignId)
	if err != nil {
		return err
	}
	return cxt.JSON(*joinCampaignCode)
}

func (c *CampaignsController) JoinCampaign(cxt *fiber.Ctx) error {
	input := types.JoinCampaignInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	token := cxt.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	playerId := token.RegisteredClaims.Subject
	campaign, err := c.CampaignsService.JoinCampaign(input, playerId)
	if err != nil {
		return err
	}
	return cxt.JSON(*campaign)
}

func NewCampaignsController(dynamodb utils.DynamoDB) CampaignsController {
	campaignsService := services.NewCampaignsService(dynamodb)

	return CampaignsController{
		CampaignsService: campaignsService,
	}
}
