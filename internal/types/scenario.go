package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

type ScenarioCreateInput struct {
	Name          *string `dynamodbav:"name,omitempty" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
	UpdatedAt     *string `dynamodbav:"updated_at" json:"updatedAt"`
}

type ScenarioPatchInput struct {
	Name          *string `dynamodbav:"name,omitempty" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
	UpdatedAt     *string `dynamodbav:"updated_at" json:"updatedAt"`
}

type ScenarioItem struct {
	Item          `dynamodbav:",inline"`
	Name          *string `dynamodbav:"name,omitempty" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
	UpdatedAt     *string `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewScenarioItem(input ScenarioCreateInput, campaignId string, scenarioId string) ScenarioItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return ScenarioItem{
		Item: Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
			Entity: constants.SCENARIO + constants.SEPERATOR + scenarioId,
		},
		Name:          input.Name,
		Groups:        input.Groups,
		ScenarioLevel: input.ScenarioLevel,
		UpdatedAt:     &updatedAt,
	}
}
