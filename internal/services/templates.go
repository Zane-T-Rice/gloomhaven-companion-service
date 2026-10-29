package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"log"
	"strings"

	"github.com/google/uuid"
)

type TemplatesService struct {
	DynamoDB utils.DynamoDB
}

// Everything the player can use: the system's templates, their own, and those of the other players in
// their campaigns.
func (s *TemplatesService) List(playerId string) ([]dto.Template, error) {
	owners, err := visibleOwners(s.DynamoDB, playerId)
	if err != nil {
		return nil, err
	}
	templates := []dto.Template{}
	for _, owner := range owners {
		templateItems := []types.TemplateItem{}
		if err := s.DynamoDB.Query(
			constants.PARENT,
			owner,
			constants.ENTITY,
			constants.TEMPLATE,
			nil,
			&templateItems,
		); err != nil {
			return nil, err
		}
		for _, templateItem := range templateItems {
			templates = append(templates, dto.NewTemplate(templateItem, playerId))
		}
	}
	return templates, nil
}

func validate(input types.TemplateCreateInput) error {
	if input.Type == nil || strings.TrimSpace(*input.Type) == "" {
		message := "A template needs a type."
		return errors.NewBadRequestError(&message)
	}
	return nil
}

// A system template is made by an admin. Otherwise the template belongs to the player.
func (s *TemplatesService) Create(input types.TemplateCreateInput, playerId string, system bool) (*dto.Template, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	owner := playerOwner(playerId)
	if system {
		owner = constants.TEMPLATE
	}
	if input.EditionId == nil || *input.EditionId == "" {
		message := "A template needs an edition."
		return nil, errors.NewBadRequestError(&message)
	}
	editions := NewEditionsService(s.DynamoDB)
	if err := editions.checkUsable(*input.EditionId, owner); err != nil {
		return nil, err
	}
	templateItem := types.NewTemplateItem(input, uuid.New().String(), owner)
	if err := s.DynamoDB.PutItem(templateItem); err != nil {
		return nil, err
	}
	template := dto.NewTemplate(templateItem, playerId)
	return &template, nil
}

// Finds a template that the player is allowed to change: their own, or a system template for an admin.
func (s *TemplatesService) findChangeable(templateId string, playerId string, admin bool) (*types.TemplateItem, error) {
	entity := constants.TEMPLATE + constants.SEPERATOR + templateId
	item := types.TemplateItem{}
	if err := s.DynamoDB.GetItem(constants.PARENT, playerOwner(playerId), constants.ENTITY, entity, &item); err == nil {
		return &item, nil
	}
	if err := s.DynamoDB.GetItem(constants.PARENT, constants.TEMPLATE, constants.ENTITY, entity, &item); err != nil {
		return nil, errors.NewNotFoundError()
	}
	if !admin {
		return nil, errors.NewForbiddenError()
	}
	return &item, nil
}

func (s *TemplatesService) Patch(input types.TemplatePatchInput, templateId string, playerId string, admin bool) (*dto.Template, error) {
	found, err := s.findChangeable(templateId, playerId, admin)
	if err != nil {
		return nil, err
	}
	if input.Type != nil {
		if err := validate(input.TemplateCreateInput); err != nil {
			return nil, err
		}
	}
	if input.EditionId != nil {
		editions := NewEditionsService(s.DynamoDB)
		if err := editions.checkUsable(*input.EditionId, found.Parent); err != nil {
			return nil, err
		}
	}
	templateItem := types.TemplateItem{}
	err = s.DynamoDB.UpdateItem(
		constants.PARENT,
		found.Parent,
		constants.ENTITY,
		found.Entity,
		input,
		&templateItem,
	)
	if err != nil {
		log.Printf("Patch failed, here's why %v", err)
		return nil, err
	}
	template := dto.NewTemplate(templateItem, playerId)
	return &template, nil
}

func (s *TemplatesService) Delete(templateId string, playerId string, admin bool) (*dto.Template, error) {
	found, err := s.findChangeable(templateId, playerId, admin)
	if err != nil {
		return nil, err
	}
	templateItem := types.TemplateItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		found.Parent,
		constants.ENTITY,
		found.Entity,
		&templateItem,
	); err != nil {
		return nil, err
	}

	template := dto.NewTemplate(templateItem, playerId)

	return &template, nil
}

func NewTemplatesService(dynamodb utils.DynamoDB) TemplatesService {
	return TemplatesService{DynamoDB: dynamodb}
}
