package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"log"
	"strings"

	"github.com/google/uuid"
)

type ScenariosService struct {
	DynamoDB utils.DynamoDB
}

func (s *ScenariosService) List(campaignId string) ([]dto.Scenario, error) {
	scenarioItems := []types.ScenarioItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.SCENARIO,
		nil,
		&scenarioItems,
	); err != nil {
		return nil, err
	}
	scenarios := []dto.Scenario{}
	for _, scenarioItem := range scenarioItems {
		scenarios = append(scenarios, dto.NewScenario(scenarioItem))
	}
	return scenarios, nil
}

func (s *ScenariosService) Create(input types.ScenarioCreateInput, campaignId string) (*dto.Scenario, error) {
	scenarioId := uuid.New().String()
	scenarioItem := types.NewScenarioItem(input, campaignId, scenarioId)
	if err := s.DynamoDB.PutItem(scenarioItem); err != nil {
		return nil, err
	}

	scenario := dto.NewScenario(scenarioItem)
	return &scenario, nil
}

func (s *ScenariosService) Patch(input types.ScenarioPatchInput, campaignId string, scenarioId string) (*dto.Scenario, error) {
	scenarioItem := types.ScenarioItem{}
	err := s.DynamoDB.UpdateItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.SCENARIO+constants.SEPERATOR+scenarioId,
		input,
		&scenarioItem,
	)
	if err != nil {
		log.Printf("Patch failed, here's why %v", err)
		return nil, err
	}
	scenario := dto.NewScenario(scenarioItem)
	return &scenario, nil
}

func (s *ScenariosService) Delete(campaignId string, scenarioId string) (*dto.Scenario, error) {
	figures := []types.FigureItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId+constants.SCENARIO+constants.SEPERATOR+scenarioId,
		constants.ENTITY,
		constants.FIGURE,
		nil,
		&figures,
	); err != nil {
		return nil, err
	}

	figuresService := NewFiguresService(s.DynamoDB)

	for _, figure := range figures {
		figuresService.Delete(campaignId, scenarioId, strings.Split(figure.Entity, constants.SEPERATOR)[2])
	}

	scenarioItem := types.ScenarioItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.SCENARIO+constants.SEPERATOR+scenarioId,
		&scenarioItem,
	); err != nil {
		return nil, err
	}

	scenario := dto.NewScenario(scenarioItem)

	return &scenario, nil
}

func NewScenariosService(dynamodb utils.DynamoDB) ScenariosService {
	return ScenariosService{DynamoDB: dynamodb}
}
