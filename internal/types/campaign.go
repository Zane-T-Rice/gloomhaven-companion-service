package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

type CampaignCreateInput struct {
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}

type CampaignPatchInput struct {
	Name      *string `dynamodbav:"name,omitempty" json:"name"`
	UpdatedAt *string `dynamodbav:"updated_at" json:"updatedAt"`
}

type CampaignItem struct {
	Item      `dynamodbav:",inline"`
	Name      *string `dynamodbav:"name,omitempty" json:"name"`
	UpdatedAt *string `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewCampaignItem(input CampaignCreateInput, campaignId string) CampaignItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return CampaignItem{
		Item: Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
			Entity: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
		},
		Name:      input.Name,
		UpdatedAt: &updatedAt,
	}
}
