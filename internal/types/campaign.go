package types

type CampaignCreateInput struct {
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}

type CampaignPatchInput struct {
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}

type CampaignItem struct {
	Item `dynamodbav:",inline"`
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}
