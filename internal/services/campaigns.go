package services

import (
	"gloomhaven-companion-service/internal/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type campaignsService struct {
	List func() ([]types.Item, error)
}

func CampaignsService(dynamodbClient *dynamodb.Client) campaignsService {
	return campaignsService{
		List: func() ([]types.Item, error) {
			return []types.Item{}, nil
		},
	}
}
