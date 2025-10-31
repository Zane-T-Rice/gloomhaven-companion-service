package services

import (
	"gloomhaven-companion-service/internal/dto"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CampaignsService struct {
	DynamoDBClient *dynamodb.Client
}

func (s CampaignsService) List() ([]dto.Campaign, error) {
	return []dto.Campaign{}, nil
}

func NewCampaignsService(dynamodbClient *dynamodb.Client) CampaignsService {
	return CampaignsService{DynamoDBClient: dynamodbClient}
}
