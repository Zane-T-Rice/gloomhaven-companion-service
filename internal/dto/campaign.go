package dto

import "gloomhaven-companion-service/internal/types"

type Campaign struct {
	Parent    string  `dynamodbav:"parent" json:"parent"`
	Entity    string  `dynamodbav:"entity" json:"entity"`
	Name      *string `dynamodbav:"name" json:"name"`
	UpdatedAt *string `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewCampaign(campaignItem types.CampaignItem) Campaign {
	return Campaign{
		Parent:    campaignItem.Parent,
		Entity:    campaignItem.Entity,
		Name:      campaignItem.Name,
		UpdatedAt: campaignItem.UpdatedAt,
	}
}
