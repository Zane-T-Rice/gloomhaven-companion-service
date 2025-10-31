package types

// type PlayerCreateInput struct {
// 	CampaignIds []string `json:"campaign_ids"`
// }

type PlayerItem struct {
	Item        `dynamodbav:",inline"`
	CampaignIds []string `dynamodbav:"campaign_ids"`
}
