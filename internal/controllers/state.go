package controllers

import (
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type StateController struct {
	StateService services.StateService
}

func (c *StateController) List(cxt *fiber.Ctx) error {
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	states, err := c.StateService.List(campaignId, scenarioId)
	if err != nil {
		return err
	}
	return cxt.JSON(states)
}

func (c *StateController) Put(cxt *fiber.Ctx) error {
	input := types.StatePutInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	if input.Key == nil || *input.Key == "" || input.Value == nil {
		return errors.NewBadRequestError(nil)
	}
	campaignId := cxt.Params("campaignId")
	scenarioId := cxt.Params("scenarioId")
	state, err := c.StateService.Put(input, campaignId, scenarioId)
	if err != nil {
		return err
	}
	return cxt.JSON(*state)
}

func NewStateController(dynamodb utils.DynamoDB) StateController {
	return StateController{
		StateService: services.NewStateService(dynamodb),
	}
}
