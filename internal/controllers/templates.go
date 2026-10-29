package controllers

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/middlewares"
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TemplatesController struct {
	TemplatesService services.TemplatesService
	EditionsService  services.EditionsService
}

func (c *TemplatesController) List(cxt *fiber.Ctx) error {
	playerId, _ := middlewares.Caller(cxt)
	templates, err := c.TemplatesService.List(playerId)
	if err != nil {
		return err
	}
	return cxt.JSON(templates)
}

// Only an admin can ask for a system template with ?scope=system.
func wantsSystem(cxt *fiber.Ctx, admin bool) (bool, error) {
	if cxt.Query("scope") != constants.SCOPE_NAME_SYSTEM {
		return false, nil
	}
	if !admin {
		return false, errors.NewForbiddenError()
	}
	return true, nil
}

func (c *TemplatesController) Create(cxt *fiber.Ctx) error {
	input := types.TemplateCreateInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	playerId, admin := middlewares.Caller(cxt)
	system, err := wantsSystem(cxt, admin)
	if err != nil {
		return err
	}
	template, err := c.TemplatesService.Create(input, playerId, system)
	if err != nil {
		return err
	}
	return cxt.JSON(*template)
}

func (c *TemplatesController) Patch(cxt *fiber.Ctx) error {
	input := types.TemplatePatchInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	playerId, admin := middlewares.Caller(cxt)
	templateId := cxt.Params("templateId")
	template, err := c.TemplatesService.Patch(input, templateId, playerId, admin)
	if err != nil {
		if strings.Contains(string(err.Error()), "The conditional request failed") {
			message := constants.UPDATED_AT_ERROR
			return errors.NewBadRequestError(&message)
		}
		return err
	}
	return cxt.JSON(*template)
}

func (c *TemplatesController) Delete(cxt *fiber.Ctx) error {
	playerId, admin := middlewares.Caller(cxt)
	templateId := cxt.Params("templateId")
	template, err := c.TemplatesService.Delete(templateId, playerId, admin)
	if err != nil {
		return err
	}
	return cxt.JSON(*template)
}

func (c *TemplatesController) ListEditions(cxt *fiber.Ctx) error {
	playerId, _ := middlewares.Caller(cxt)
	editions, err := c.EditionsService.List(playerId)
	if err != nil {
		return err
	}
	return cxt.JSON(editions)
}

func (c *TemplatesController) CreateEdition(cxt *fiber.Ctx) error {
	input := types.EditionCreateInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	playerId, admin := middlewares.Caller(cxt)
	system, err := wantsSystem(cxt, admin)
	if err != nil {
		return err
	}
	edition, err := c.EditionsService.Create(input, playerId, system)
	if err != nil {
		return err
	}
	return cxt.JSON(*edition)
}

func (c *TemplatesController) PatchEdition(cxt *fiber.Ctx) error {
	input := types.EditionPatchInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	playerId, admin := middlewares.Caller(cxt)
	edition, err := c.EditionsService.Patch(input, cxt.Params("editionId"), playerId, admin)
	if err != nil {
		if strings.Contains(string(err.Error()), "The conditional request failed") {
			message := constants.UPDATED_AT_ERROR
			return errors.NewBadRequestError(&message)
		}
		return err
	}
	return cxt.JSON(*edition)
}

func (c *TemplatesController) DeleteEdition(cxt *fiber.Ctx) error {
	playerId, admin := middlewares.Caller(cxt)
	edition, err := c.EditionsService.Delete(cxt.Params("editionId"), playerId, admin)
	if err != nil {
		return err
	}
	return cxt.JSON(*edition)
}

func NewTemplatesController(dynamodb utils.DynamoDB) TemplatesController {
	return TemplatesController{
		TemplatesService: services.NewTemplatesService(dynamodb),
		EditionsService:  services.NewEditionsService(dynamodb),
	}
}
