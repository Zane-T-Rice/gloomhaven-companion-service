package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CampaignsService struct {
	DynamoDB utils.DynamoDB
}

func (s *CampaignsService) List(playerId string) ([]dto.Campaign, error) {
	playerCampaigns := []types.PlayerCampaignItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.PLAYER+constants.SEPERATOR+playerId,
		constants.ENTITY,
		constants.CAMPAIGN,
		nil,
		&playerCampaigns,
	); err != nil {
		return nil, err
	}
	campaigns := []dto.Campaign{}
	for _, playerCampaign := range playerCampaigns {
		item := types.CampaignItem{}
		s.DynamoDB.GetItem(
			constants.PARENT,
			playerCampaign.Entity,
			constants.ENTITY,
			playerCampaign.Entity,
			&item,
		)
		campaigns = append(campaigns, dto.NewCampaign(item))
	}
	return campaigns, nil
}

func (s *CampaignsService) Create(input types.CampaignCreateInput, playerId string) (*dto.Campaign, error) {
	campaignId := uuid.New().String()
	campaignItem := types.CampaignItem{
		Item: types.Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
		},
		Name: input.Name,
	}
	if err := s.DynamoDB.PutItem(campaignItem); err != nil {
		return nil, err
	}
	playerCampaignItem := types.PlayerCampaignItem{
		Item: types.Item{
			Parent: constants.PLAYER + constants.SEPERATOR + playerId,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
		},
	}
	if err := s.DynamoDB.PutItem(playerCampaignItem); err != nil {
		return nil, err
	}

	campaign := dto.NewCampaign(campaignItem)
	return &campaign, nil
}

func (s *CampaignsService) Patch(input types.CampaignPatchInput, campaignId string) (*dto.Campaign, error) {
	campaignItem := types.CampaignItem{}
	s.DynamoDB.UpdateItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		input,
		&campaignItem,
	)
	campaign := dto.NewCampaign(campaignItem)
	return &campaign, nil
}

func (s *CampaignsService) Delete(campaignId string) (*dto.Campaign, error) {
	scenarios := []types.ScenarioItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.SCENARIO,
		nil,
		&scenarios,
	); err != nil {
		return nil, err
	}

	scenariosService := NewScenariosService(s.DynamoDB)

	for _, scenario := range scenarios {
		scenariosService.Delete(campaignId, strings.Split(scenario.Entity, constants.SEPERATOR)[2])
	}

	playerCampaigns := []types.PlayerCampaignItem{}
	if err := s.DynamoDB.Query(
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.PARENT,
		constants.PLAYER,
		&constants.ENTITY_INDEX,
		&playerCampaigns,
	); err != nil {
		return nil, err
	}

	for _, playerCampaign := range playerCampaigns {
		if err := s.DynamoDB.DeleteItem(
			constants.PARENT,
			playerCampaign.Parent,
			constants.ENTITY,
			playerCampaign.Entity,
			&types.PlayerCampaignItem{},
		); err != nil {
			return nil, err
		}
	}

	joinCampaigns := []types.JoinCampaignItem{}
	if err := s.DynamoDB.Query(
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.PARENT,
		constants.JOIN,
		&constants.ENTITY_INDEX,
		&joinCampaigns,
	); err != nil {
		return nil, err
	}

	for _, joinCampaign := range joinCampaigns {
		if err := s.DynamoDB.DeleteItem(
			constants.PARENT,
			joinCampaign.Parent,
			constants.ENTITY,
			joinCampaign.Entity,
			&types.JoinCampaignItem{},
		); err != nil {
			return nil, err
		}
	}

	campaignItem := types.CampaignItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		&campaignItem,
	); err != nil {
		return nil, err
	}

	campaign := dto.NewCampaign(campaignItem)

	return &campaign, nil
}

func (s *CampaignsService) JoinCampaign(input types.JoinCampaignInput, playerId string) (*dto.Campaign, error) {
	now := time.Now().Unix()
	joinCampaignItems := []types.JoinCampaignItem{}
	if err := s.DynamoDB.Query(
		constants.PARENT,
		constants.JOIN+constants.SEPERATOR+*input.Code,
		constants.ENTITY,
		constants.CAMPAIGN,
		nil,
		&joinCampaignItems,
	); err != nil {
		return nil, err
	}

	if now-*joinCampaignItems[0].CreatedAt > 300 {
		return nil, errors.NewForbiddenError()
	}

	campaignId := strings.Split(joinCampaignItems[0].Entity, "#")[2]
	playerCampaignItem := types.PlayerCampaignItem{
		Item: types.Item{
			Parent: constants.PLAYER + constants.SEPERATOR + playerId,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
		},
	}
	if err := s.DynamoDB.PutItem(playerCampaignItem); err != nil {
		return nil, err
	}

	campaignItem := types.CampaignItem{}
	if err := s.DynamoDB.GetItem(
		constants.PARENT,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		&campaignItem,
	); err != nil {
		return nil, err
	}

	campaign := dto.NewCampaign(campaignItem)
	return &campaign, nil
}

func (s *CampaignsService) CreateJoinCode(campaignId string) (*dto.JoinCampaign, error) {
	now := time.Now().Unix()
	joinCode := utils.GenerateRandomString(10)
	joinCampaignItem := types.JoinCampaignItem{
		Item: types.Item{
			Parent: constants.JOIN + constants.SEPERATOR + joinCode,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
		},
		JoinCampaignInput: types.JoinCampaignInput{
			Code: &joinCode,
		},
		CreatedAt: &now,
	}
	if err := s.DynamoDB.PutItem(joinCampaignItem); err != nil {
		return nil, err
	}
	joinCampaign := dto.NewJoinCampaign(joinCampaignItem)
	return &joinCampaign, nil
}

func NewCampaignsService(dynamodb utils.DynamoDB) CampaignsService {
	return CampaignsService{DynamoDB: dynamodb}
}
