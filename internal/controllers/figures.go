package controllers

import (
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type FiguresController struct {
	FiguresService services.FiguresService
}

func (c *FiguresController) List(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	figures, err := c.FiguresService.List(campaignId, scenarioId)
	if err != nil {
		return err
	}
	return cxt.JSON(figures)
}

func (c *FiguresController) Create(cxt *fiber.Ctx) error {
	input := types.FigureCreateInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	if input.Damage == nil || input.MaximumHP == nil || input.Class == nil {
		return errors.NewBadRequestError()
	}
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	figure, err := c.FiguresService.Create(input, campaignId, scenarioId)
	if err != nil {
		return err
	}
	return cxt.JSON(*figure)
}

func (c *FiguresController) Patch(cxt *fiber.Ctx) error {
	input := types.FigurePatchInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	if input.Damage == nil && input.MaximumHP == nil && input.Name == nil {
		return errors.NewBadRequestError()
	}
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	figureId := cxt.Params("figureId")
	scenario, err := c.FiguresService.Patch(input, campaignId, scenarioId, figureId)
	if err != nil {
		return err
	}
	return cxt.JSON(*scenario)
}

func (c *FiguresController) Delete(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	figureId := cxt.Params("figureId")
	scenario, err := c.FiguresService.Delete(campaignId, scenarioId, figureId)
	if err != nil {
		return err
	}
	return cxt.JSON(*scenario)
}

func NewFiguresController(dynamodb utils.DynamoDB) FiguresController {
	scenariosService := services.NewFiguresService(dynamodb)

	return FiguresController{
		FiguresService: scenariosService,
	}
}
