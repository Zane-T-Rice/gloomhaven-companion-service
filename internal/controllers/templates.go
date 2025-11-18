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
	templateId := cxt.Params("templateId")
	template, err := c.TemplatesService.Patch(input, templateId)
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
	templateId := cxt.Params("templateId")
	template, err := c.TemplatesService.Delete(templateId)
	if err != nil {
		return err
	}
	return cxt.JSON(*template)
}

func NewTemplatesController(dynamodb utils.DynamoDB) TemplatesController {
	templatesService := services.NewTemplatesService(dynamodb)

	return TemplatesController{
		TemplatesService: templatesService,
	}
}
