package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

	"github.com/google/uuid"
)

type FiguresService struct {
	DynamoDB utils.DynamoDB
}

func (s *FiguresService) List(campaignId string, scenarioId string) ([]dto.Figure, error) {
	figureItems := []types.FigureItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId+constants.SCENARIO+constants.SEPERATOR+scenarioId,
		constants.ENTITY,
		constants.FIGURE,
		nil,
		&figureItems,
	); err != nil {
		return nil, err
	}
	figures := []dto.Figure{}
	for _, figureItem := range figureItems {
		figures = append(figures, dto.Figure{
			Parent:    figureItem.Parent,
			Entity:    figureItem.Entity,
			Name:      figureItem.Name,
			MaximumHP: figureItem.MaximumHP,
			Damage:    figureItem.Damage,
		})
	}
	return figures, nil
}

func (s *FiguresService) Create(input types.FigureCreateInput, campaignId string, scenarioId string) (*dto.Figure, error) {
	figureId := uuid.New().String()
	figureItem := types.FigureItem{
		Item: types.Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId + constants.SCENARIO + constants.SEPERATOR + scenarioId,
			Entity: constants.FIGURE + constants.SEPERATOR + figureId,
		},
		Name:      input.Name,
		MaximumHP: input.MaximumHP,
		Damage:    input.Damage,
	}
	if err := s.DynamoDB.PutItem(figureItem); err != nil {
		return nil, err
	}

	figure := dto.Figure{
		Parent:    figureItem.Parent,
		Entity:    figureItem.Entity,
		Name:      figureItem.Name,
		MaximumHP: figureItem.MaximumHP,
		Damage:    figureItem.Damage,
	}
	return &figure, nil
}

func (s *FiguresService) Patch(input types.FigurePatchInput, campaignId string, scenarioId string, figureId string) (*dto.Figure, error) {
	figureItem := types.FigureItem{}
	s.DynamoDB.UpdateItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId+constants.SCENARIO+constants.SEPERATOR+scenarioId,
		constants.ENTITY,
		constants.FIGURE+constants.SEPERATOR+figureId,
		input,
		&figureItem,
	)
	figure := dto.Figure{
		Parent:    figureItem.Parent,
		Entity:    figureItem.Entity,
		Name:      figureItem.Name,
		MaximumHP: figureItem.MaximumHP,
		Damage:    figureItem.Damage,
	}
	return &figure, nil
}

func (s *FiguresService) Delete(campaignId string, scenarioId string, figureId string) (*dto.Figure, error) {
	figureItem := types.FigureItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId+constants.SCENARIO+constants.SEPERATOR+scenarioId,
		constants.ENTITY,
		constants.FIGURE+constants.SEPERATOR+figureId,
		&figureItem,
	); err != nil {
		return nil, err
	}

	figure := dto.Figure{
		Parent: figureItem.Parent,
		Entity: figureItem.Entity,
		Name:   figureItem.Name,
	}

	return &figure, nil
}

func NewFiguresService(dynamodb utils.DynamoDB) FiguresService {
	return FiguresService{DynamoDB: dynamodb}
}
