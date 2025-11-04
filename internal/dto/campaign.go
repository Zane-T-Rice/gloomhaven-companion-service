package dto

import "gloomhaven-companion-service/internal/types"

type Campaign struct {
	Parent string  `dynamodbav:"parent" json:"parent"`
	Entity string  `dynamodbav:"entity" json:"entity"`
	Name   *string `dynamodbav:"name" json:"name"`
}

func NewCampaign(campaignItem types.CampaignItem) Campaign {
	return Campaign{
		Parent: campaignItem.Parent,
		Entity: campaignItem.Entity,
		Name:   campaignItem.Name,
	}
}
