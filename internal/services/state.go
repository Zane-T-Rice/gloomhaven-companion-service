package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
)

type StateService struct {
	DynamoDB utils.DynamoDB
}

func (s *StateService) List(campaignId string, scenarioId string) ([]dto.State, error) {
	stateItems := []types.StateItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId+constants.SCENARIO+constants.SEPERATOR+scenarioId,
		constants.ENTITY,
		constants.STATE,
		nil,
		&stateItems,
	); err != nil {
		return nil, err
	}
	states := []dto.State{}
	for _, stateItem := range stateItems {
		states = append(states, dto.NewState(stateItem))
	}
	return states, nil
}

func (s *StateService) Put(input types.StatePutInput, campaignId string, scenarioId string) (*dto.State, error) {
	stateItem := types.NewStateItem(input, campaignId, scenarioId)
	if err := s.DynamoDB.PutItem(stateItem); err != nil {
		return nil, err
	}
	state := dto.NewState(stateItem)
	return &state, nil
}

func (s *StateService) DeleteAll(campaignId string, scenarioId string) error {
	stateItems := []types.StateItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId+constants.SCENARIO+constants.SEPERATOR+scenarioId,
		constants.ENTITY,
		constants.STATE,
		nil,
		&stateItems,
	); err != nil {
		return err
	}
	for _, stateItem := range stateItems {
		deleted := types.StateItem{}
		if err := s.DynamoDB.DeleteItem(
			constants.PARENT,
			stateItem.Parent,
			constants.ENTITY,
			stateItem.Entity,
			&deleted,
		); err != nil {
			return err
		}
	}
	return nil
}

func NewStateService(dynamodb utils.DynamoDB) StateService {
	return StateService{DynamoDB: dynamodb}
}
