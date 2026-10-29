package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"strings"

	"github.com/google/uuid"
)

type EditionsService struct {
	DynamoDB utils.DynamoDB
}

func (s *EditionsService) List(playerId string) ([]dto.Edition, error) {
	owners, err := visibleOwners(s.DynamoDB, playerId)
	if err != nil {
		return nil, err
	}
	editions := []dto.Edition{dto.NewDefaultEdition()}
	for _, owner := range owners {
		// The system's editions are kept under constants.EDITION, not constants.TEMPLATE.
		parent := owner
		if owner == constants.TEMPLATE {
			parent = constants.EDITION
		}
		items := []types.EditionItem{}
		if err := s.DynamoDB.Query(
			constants.PARENT,
			parent,
			constants.ENTITY,
			constants.EDITION,
			nil,
			&items,
		); err != nil {
			return nil, err
		}
		for _, item := range items {
			editions = append(editions, dto.NewEdition(item, playerId))
		}
	}
	return editions, nil
}

// Which parent a new edition goes under.
func (s *EditionsService) ownerFor(playerId string, system bool) string {
	if system {
		return constants.EDITION
	}
	return playerOwner(playerId)
}

func validateName(name *string) error {
	if name == nil || strings.TrimSpace(*name) == "" {
		message := "An edition needs a name."
		return errors.NewBadRequestError(&message)
	}
	return nil
}

func (s *EditionsService) Create(input types.EditionCreateInput, playerId string, system bool) (*dto.Edition, error) {
	if err := validateName(input.Name); err != nil {
		return nil, err
	}
	editionItem := types.NewEditionItem(input, uuid.New().String(), s.ownerFor(playerId, system))
	if err := s.DynamoDB.PutItem(editionItem); err != nil {
		return nil, err
	}
	edition := dto.NewEdition(editionItem, playerId)
	return &edition, nil
}

// Finds an edition that the player is allowed to change.
func (s *EditionsService) findChangeable(editionId string, playerId string, admin bool) (*types.EditionItem, error) {
	if editionId == constants.DEFAULT_EDITION_ID {
		return nil, errors.NewForbiddenError()
	}
	entity := constants.EDITION + constants.SEPERATOR + editionId
	item := types.EditionItem{}
	if err := s.DynamoDB.GetItem(constants.PARENT, playerOwner(playerId), constants.ENTITY, entity, &item); err == nil {
		return &item, nil
	}
	if err := s.DynamoDB.GetItem(constants.PARENT, constants.EDITION, constants.ENTITY, entity, &item); err != nil {
		return nil, errors.NewNotFoundError()
	}
	if !admin {
		return nil, errors.NewForbiddenError()
	}
	return &item, nil
}

func (s *EditionsService) Patch(input types.EditionPatchInput, editionId string, playerId string, admin bool) (*dto.Edition, error) {
	if err := validateName(input.Name); err != nil {
		return nil, err
	}
	found, err := s.findChangeable(editionId, playerId, admin)
	if err != nil {
		return nil, err
	}
	editionItem := types.EditionItem{}
	if err := s.DynamoDB.UpdateItem(
		constants.PARENT,
		found.Parent,
		constants.ENTITY,
		found.Entity,
		input,
		&editionItem,
	); err != nil {
		return nil, err
	}
	edition := dto.NewEdition(editionItem, playerId)
	return &edition, nil
}

// An edition that still has templates in it is not deleted, so that no template is left without one.
func (s *EditionsService) Delete(editionId string, playerId string, admin bool) (*dto.Edition, error) {
	found, err := s.findChangeable(editionId, playerId, admin)
	if err != nil {
		return nil, err
	}
	templateOwner := found.Parent
	if templateOwner == constants.EDITION {
		templateOwner = constants.TEMPLATE
	}
	templates := []types.TemplateItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		templateOwner,
		constants.ENTITY,
		constants.TEMPLATE,
		nil,
		&templates,
	); err != nil {
		return nil, err
	}
	for _, template := range templates {
		if template.EditionId != nil && *template.EditionId == editionId {
			message := "This edition still has templates in it. Move or delete them first."
			return nil, errors.NewBadRequestError(&message)
		}
	}
	editionItem := types.EditionItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		found.Parent,
		constants.ENTITY,
		found.Entity,
		&editionItem,
	); err != nil {
		return nil, err
	}
	edition := dto.NewEdition(editionItem, playerId)
	return &edition, nil
}

// Whether templates belonging to the owner can be put in the edition. That is the case for the default
// edition, for the system's editions, and for the owner's own editions.
func (s *EditionsService) checkUsable(editionId string, owner string) error {
	if editionId == constants.DEFAULT_EDITION_ID {
		return nil
	}
	entity := constants.EDITION + constants.SEPERATOR + editionId
	item := types.EditionItem{}
	if err := s.DynamoDB.GetItem(constants.PARENT, constants.EDITION, constants.ENTITY, entity, &item); err == nil {
		return nil
	}
	if owner != constants.TEMPLATE {
		if err := s.DynamoDB.GetItem(constants.PARENT, owner, constants.ENTITY, entity, &item); err == nil {
			return nil
		}
	}
	message := "That edition does not exist."
	return errors.NewBadRequestError(&message)
}

func NewEditionsService(dynamodb utils.DynamoDB) EditionsService {
	return EditionsService{DynamoDB: dynamodb}
}
