package controllers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ScenariosController struct {
	ScenariosService services.ScenariosService
}

func (c *ScenariosController) List(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	scenarios, err := c.ScenariosService.List(campaignId)
	if err != nil {
		return err
	}
	return cxt.JSON(scenarios)
}

func (c *ScenariosController) Create(cxt *fiber.Ctx) error {
	input := types.ScenarioCreateInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	if input.Name == nil {
		return errors.NewBadRequestError(nil)
	}
	campaignId := cxt.Params("campaignId")
	scenario, err := c.ScenariosService.Create(input, campaignId)
	if err != nil {
		return err
	}
	return cxt.JSON(*scenario)
}

func (c *ScenariosController) Patch(cxt *fiber.Ctx) error {
	input := types.ScenarioPatchInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	if input.Name == nil {
		return errors.NewBadRequestError(nil)
	}
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	scenario, err := c.ScenariosService.Patch(input, campaignId, scenarioId)
	if err != nil {
		if strings.Contains(string(err.Error()), "The conditional request failed") {
			message := constants.UPDATED_AT_ERROR
			return errors.NewBadRequestError(&message)
		}
		return err
	}
	return cxt.JSON(*scenario)
}

func (c *ScenariosController) Delete(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	scenario, err := c.ScenariosService.Delete(campaignId, scenarioId)
	if err != nil {
		return err
	}
	return cxt.JSON(*scenario)
}

func NewScenariosController(dynamodb utils.DynamoDB) ScenariosController {
	scenariosService := services.NewScenariosService(dynamodb)

	return ScenariosController{
		ScenariosService: scenariosService,
	}
}
