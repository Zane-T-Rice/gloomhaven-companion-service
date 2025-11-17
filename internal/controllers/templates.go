package controllers

import (
	"gloomhaven-companion-service/internal/services"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type TemplatesController struct {
	TemplatesService services.TemplatesService
}

func (c *TemplatesController) List(cxt *fiber.Ctx) error {
	templates, err := c.TemplatesService.List()
	if err != nil {
		return err
	}
	return cxt.JSON(templates)
}

func (c *TemplatesController) Create(cxt *fiber.Ctx) error {
	input := types.TemplateCreateInput{}
	if err := cxt.BodyParser(&input); err != nil {
		return err
	}
	// if input.Damage == nil || input.MaximumHP == nil || input.Name == nil {
	// 	return errors.NewBadRequestError()
	// }
	template, err := c.TemplatesService.Create(input)
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
	// if input.Damage == nil && input.MaximumHP == nil && input.Name == nil {
	// 	return errors.NewBadRequestError()
	// }
	templateId := cxt.Params("templateId")
	scenario, err := c.TemplatesService.Patch(input, templateId)
	if err != nil {
		return err
	}
	return cxt.JSON(*scenario)
}

func (c *TemplatesController) Delete(cxt *fiber.Ctx) error {
	templateId := cxt.Params("templateId")
	scenario, err := c.TemplatesService.Delete(templateId)
	if err != nil {
		return err
	}
	return cxt.JSON(*scenario)
}

func NewTemplatesController(dynamodb utils.DynamoDB) TemplatesController {
	scenariosService := services.NewTemplatesService(dynamodb)

	return TemplatesController{
		TemplatesService: scenariosService,
	}
}
