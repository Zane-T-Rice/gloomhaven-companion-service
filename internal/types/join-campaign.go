package types

type JoinCampaignInput struct {
	Code *string `dynamodbav:"code" json:"code"`
}

type JoinCampaignItem struct {
	Item              `dynamodbav:",inline"`
	JoinCampaignInput `dynamodbav:",inline"`
	CreatedAt         *int64 `dynamodbav:"code" json:"code"`
}
