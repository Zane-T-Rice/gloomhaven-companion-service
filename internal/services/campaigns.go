package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
	"strings"

	"github.com/google/uuid"
)

type CampaignsService struct {
	DynamoDB utils.DynamoDB
}

func (s *CampaignsService) List() ([]dto.Campaign, error) {
	return []dto.Campaign{}, nil
}

func (s *CampaignsService) Create(input types.CampaignCreateInput, playerId string) (*dto.Campaign, error) {
	campaignId := uuid.New().String()
	item := types.CampaignItem{
		Item: types.Item{
			Parent: constants.ROOT,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
		},
		Name: input.Name,
	}
	campaign := dto.Campaign{
		Parent: item.Parent,
		Entity: item.Entity,
		Name:   item.Name,
	}
	if err := s.DynamoDB.PutItem(item); err != nil {
		return nil, err
	}

	// Also add the campaign to the list of campaigns under Player
	player := types.PlayerItem{}
	if err := s.DynamoDB.GetItem(
		constants.PARENT,
		constants.ROOT,
		constants.PLAYER,
		constants.PLAYER+constants.SEPERATOR+playerId,
		&player,
	); err != nil {
		if !strings.Contains(err.Error(), constants.NOT_FOUND_ERROR_MESSAGE) {
			return nil, err
		}

		player := types.PlayerItem{
			Item: types.Item{
				Parent: constants.ROOT,
				Entity: constants.PLAYER + constants.SEPERATOR + playerId,
			},
			CampaignIds: []string{campaignId},
		}
		if err := s.DynamoDB.PutItem(player); err != nil {
			return nil, err
		}
	} else {
		player.CampaignIds = append(player.CampaignIds, campaignId)
		if err := s.DynamoDB.PutItem(player); err != nil {
			return nil, err
		}
	}

	return &campaign, nil
}

func NewCampaignsService(dynamodb utils.DynamoDB) CampaignsService {
	return CampaignsService{DynamoDB: dynamodb}
}
