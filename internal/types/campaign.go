package types

type CampaignCreateInput struct {
	Name string `dynamodbav:"name"`
}

type CampaignPatchInput struct {
	Name string `dynamodbav:"name"`
}

type CampaignItem struct {
	Item `dynamodbav:",inline"`
	Name string `dynamodbav:"name"`
}
