package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"log"

	"github.com/google/uuid"
)

type ScenariosService struct {
	DynamoDB utils.DynamoDB
}

func (s *ScenariosService) List(campaignId string) ([]dto.Scenario, error) {
	log.Printf("Looking up scenarios for campaign %s", campaignId)
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
		scenarios = append(scenarios, dto.Scenario{
			Parent: scenarioItem.Parent,
			Entity: scenarioItem.Entity,
			Name:   scenarioItem.Name,
		})
	}
	return scenarios, nil
}

func (s *ScenariosService) Create(input types.ScenarioCreateInput, campaignId string) (*dto.Scenario, error) {
	scenarioId := uuid.New().String()
	scenarioItem := types.ScenarioItem{
		Item: types.Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
			Entity: constants.SCENARIO + constants.SEPERATOR + scenarioId,
		},
		Name: input.Name,
	}
	if err := s.DynamoDB.PutItem(scenarioItem); err != nil {
		return nil, err
	}

	scenario := dto.Scenario{
		Parent: scenarioItem.Parent,
		Entity: scenarioItem.Entity,
		Name:   scenarioItem.Name,
	}
	return &scenario, nil
}

func (s *ScenariosService) Patch(input types.ScenarioPatchInput, campaignId string, scenarioId string) (*dto.Scenario, error) {
	scenarioItem := types.ScenarioItem{}
	s.DynamoDB.UpdateItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.SCENARIO+constants.SEPERATOR+scenarioId,
		input,
		&scenarioItem,
	)
	scenario := dto.Scenario{
		Parent: scenarioItem.Parent,
		Entity: scenarioItem.Entity,
		Name:   scenarioItem.Name,
	}
	return &scenario, nil
}

func (s *ScenariosService) Delete(campaignId string, scenarioId string) (*dto.Scenario, error) {
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

	scenario := dto.Scenario{
		Parent: scenarioItem.Parent,
		Entity: scenarioItem.Entity,
		Name:   scenarioItem.Name,
	}

	return &scenario, nil
}

func NewScenariosService(dynamodb utils.DynamoDB) ScenariosService {
	return ScenariosService{DynamoDB: dynamodb}
}
