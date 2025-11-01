package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

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
			constants.ROOT,
			constants.ENTITY,
			playerCampaign.Entity,
			&item,
		)
		campaigns = append(campaigns, dto.Campaign{
			Parent: item.Parent,
			Entity: item.Entity,
			Name:   item.Name,
		})
	}
	return campaigns, nil
}

func (s *CampaignsService) Create(input types.CampaignCreateInput, playerId string) (*dto.Campaign, error) {
	campaignId := uuid.New().String()
	campaignItem := types.CampaignItem{
		Item: types.Item{
			Parent: constants.ROOT,
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

	campaign := dto.Campaign{
		Parent: campaignItem.Parent,
		Entity: campaignItem.Entity,
		Name:   campaignItem.Name,
	}
	return &campaign, nil
}

func (s *CampaignsService) Patch(input types.CampaignPatchInput, campaignId string) (*dto.Campaign, error) {
	campaignItem := types.CampaignItem{}
	s.DynamoDB.UpdateItem(
		constants.PARENT,
		constants.ROOT,
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		input,
		&campaignItem,
	)
	campaign := dto.Campaign{
		Parent: campaignItem.Parent,
		Entity: campaignItem.Entity,
		Name:   campaignItem.Name,
	}
	return &campaign, nil
}

func (s *CampaignsService) Delete(campaignId string) (*dto.Campaign, error) {
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

	campaignItem := types.CampaignItem{}
	if err := s.DynamoDB.DeleteItem(
		constants.PARENT,
		constants.ROOT,
		constants.ENTITY,
		constants.CAMPAIGN+constants.SEPERATOR+campaignId,
		&campaignItem,
	); err != nil {
		return nil, err
	}

	campaign := dto.Campaign{
		Parent: campaignItem.Parent,
		Entity: campaignItem.Entity,
		Name:   campaignItem.Name,
	}

	return &campaign, nil
}

func NewCampaignsService(dynamodb utils.DynamoDB) CampaignsService {
	return CampaignsService{DynamoDB: dynamodb}
}
