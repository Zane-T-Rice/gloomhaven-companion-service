package services

import (
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/utils"
)

type CampaignsService struct {
	DynamoDB utils.DynamoDB
}

func (s CampaignsService) List() ([]dto.Campaign, error) {
	return []dto.Campaign{}, nil
}

func NewCampaignsService(dynamodb utils.DynamoDB) CampaignsService {
	return CampaignsService{DynamoDB: dynamodb}
}
