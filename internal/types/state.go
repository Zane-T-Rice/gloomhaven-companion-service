package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

// State is one key/value pair of table state for a scenario (the round, an element, a player's
// initiative, and so on). Every key is its own item and is overwritten whole, with no updatedAt check,
// so players changing different keys at the same time can never conflict with each other.
type StatePutInput struct {
	Key   *string `dynamodbav:"-" json:"key"`
	Value *string `dynamodbav:"value" json:"value"`
}

type StateItem struct {
	Item      `dynamodbav:",inline"`
	Value     *string `dynamodbav:"value" json:"value"`
	UpdatedAt *string `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewStateItem(input StatePutInput, campaignId string, scenarioId string) StateItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return StateItem{
		Item: Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId + constants.SCENARIO + constants.SEPERATOR + scenarioId,
			Entity: constants.STATE + constants.SEPERATOR + *input.Key,
		},
		Value:     input.Value,
		UpdatedAt: &updatedAt,
	}
}
