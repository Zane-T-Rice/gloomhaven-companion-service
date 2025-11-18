package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"log"

	"github.com/google/uuid"
)

type TemplatesService struct {
	DynamoDB utils.DynamoDB
}

func (s *TemplatesService) List() ([]dto.Template, error) {
	templateItems := []types.TemplateItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.TEMPLATE,
		constants.ENTITY,
		constants.TEMPLATE,
		nil,
		&templateItems,
	); err != nil {
		return nil, err
	}
	templates := []dto.Template{}
	for _, templateItem := range templateItems {
		templates = append(templates, dto.NewTemplate(templateItem))
	}
	return templates, nil
}

func (s *TemplatesService) Create(input types.TemplateCreateInput) (*dto.Template, error) {
	templateId := uuid.New().String()
	templateItem := types.NewTemplateItem(input, templateId)
	if err := s.DynamoDB.PutItem(templateItem); err != nil {
		return nil, err
	}
	template := dto.NewTemplate(templateItem)
	return &template, nil
}

func (s *TemplatesService) Patch(input types.TemplatePatchInput, templateId string) (*dto.Template, error) {
	templateItem := types.TemplateItem{}
	err := s.DynamoDB.UpdateItem(
		constants.PARENT,
		constants.TEMPLATE,
		constants.ENTITY,
		constants.TEMPLATE+constants.SEPERATOR+templateId,
		input,
		&templateItem,
	)
	if err != nil {
		log.Printf("Patch failed, here's why %v", err)
		return nil, err
	}
	template := dto.NewTemplate(templateItem)
	return &template, nil
}

func (s *TemplatesService) Delete(templateId string) (*dto.Template, error) {
	templateItem := types.TemplateItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		constants.TEMPLATE,
		constants.ENTITY,
		constants.TEMPLATE+constants.SEPERATOR+templateId,
		&templateItem,
	); err != nil {
		return nil, err
	}

	template := dto.NewTemplate(templateItem)

	return &template, nil
}

func NewTemplatesService(dynamodb utils.DynamoDB) TemplatesService {
	return TemplatesService{DynamoDB: dynamodb}
}
