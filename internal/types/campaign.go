package types

type CampaignCreateInput struct {
	Name string `json:"name"`
}

type CampaignItem struct {
	Item `dynamodbav:",inline"`
	Name string `dynamodbav:"name"`
}
