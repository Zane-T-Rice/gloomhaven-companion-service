package dto

import "gloomhaven-companion-service/internal/types"

type JoinCampaign struct {
	Code *string `dynamodbav:"code" json:"code"`
}

func NewJoinCampaign(joinCampaignItem types.JoinCampaignItem) JoinCampaign {
	return JoinCampaign{
		Code: joinCampaignItem.Code,
	}
}
